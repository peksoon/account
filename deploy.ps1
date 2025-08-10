# 스마트 가계부 Docker 배포 스크립트 (Windows PowerShell)
# Host Network 모드로 실행

param(
    [switch]$Clean,
    [switch]$ForceClean,
    [switch]$Stop,
    [switch]$Status,
    [switch]$Help
)

# 프로젝트 정보
$ProjectName = "iksoon-account"
$BackendImage = "$ProjectName-backend"
$FrontendImage = "$ProjectName-frontend"
$DataDir = "./data"

# 색깔 함수들
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Log-Info($message) {
    Write-ColorOutput Blue "[INFO] $message"
}

function Log-Success($message) {
    Write-ColorOutput Green "[SUCCESS] $message"
}

function Log-Warning($message) {
    Write-ColorOutput Yellow "[WARNING] $message"
}

function Log-Error($message) {
    Write-ColorOutput Red "[ERROR] $message"
}

# 함수: Docker 설치 확인
function Check-Docker {
    Log-Info "Docker 설치 상태 확인 중..."
    
    try {
        $dockerVersion = docker --version 2>$null
        if (-not $dockerVersion) {
            throw "Docker not found"
        }
    }
    catch {
        Log-Error "Docker가 설치되어 있지 않습니다."
        Log-Info "Docker 설치: https://docs.docker.com/desktop/windows/"
        exit 1
    }
    
    try {
        $composeVersion = docker-compose --version 2>$null
        if (-not $composeVersion) {
            throw "Docker Compose not found"
        }
    }
    catch {
        Log-Error "Docker Compose가 설치되어 있지 않습니다."
        Log-Info "Docker Desktop에 포함되어 있어야 합니다."
        exit 1
    }
    
    Log-Success "Docker 및 Docker Compose가 설치되어 있습니다."
}

# 함수: 기존 컨테이너 정리
function Cleanup-Containers($cleanImages = $false, $forceClean = $false) {
    Log-Info "기존 컨테이너 정리 중..."
    
    # 실행 중인 컨테이너 확인 및 정지
    $runningContainers = docker ps --format "table {{.Names}}" | Select-String $ProjectName
    if ($runningContainers) {
        Log-Warning "실행 중인 $ProjectName 컨테이너를 정지합니다."
        docker-compose down
    }
    
    # 기존 이미지 제거 (선택적)
    if ($cleanImages) {
        Log-Info "기존 이미지 제거 중..."
        try {
            docker rmi $BackendImage 2>$null
            docker rmi $FrontendImage 2>$null
            docker system prune -f
        }
        catch {
            # 이미지가 없어도 계속 진행
        }
    }
    
    # 완전한 캐시 제거 (--force-clean)
    if ($forceClean) {
        Log-Warning "🧹 모든 Docker 캐시 및 이미지 제거 중..."
        
        try {
            # 프로젝트 관련 컨테이너 강제 제거
            $containers = docker ps -a --format "table {{.ID}} {{.Names}}" | Select-String $ProjectName
            if ($containers) {
                $containerIds = $containers | ForEach-Object { ($_ -split '\s+')[0] }
                docker rm -f $containerIds 2>$null
            }
            
            # 프로젝트 관련 이미지 제거
            docker rmi $BackendImage $FrontendImage 2>$null
            
            # 사용하지 않는 모든 것 제거
            docker system prune -a -f --volumes
            
            # 빌드 캐시 완전 제거
            docker builder prune -a -f
            
            # 네트워크 정리
            docker network prune -f
            
            Log-Success "모든 Docker 캐시가 제거되었습니다."
        }
        catch {
            Log-Warning "일부 캐시 제거 중 오류가 발생했지만 계속 진행합니다."
        }
    }
    
    Log-Success "컨테이너 정리 완료"
}

# 함수: 데이터 디렉토리 생성
function Create-DataDirectory {
    Log-Info "데이터 디렉토리 생성 중..."
    
    if (-not (Test-Path $DataDir)) {
        New-Item -ItemType Directory -Path $DataDir -Force | Out-Null
    }
    
    Log-Success "데이터 디렉토리 생성 완료: $DataDir"
}

# 함수: 이미지 빌드
function Build-Images($forceClean = $false) {
    Log-Info "Docker 이미지 빌드 시작..."
    
    # 캐시 제거 옵션 확인
    $buildArgs = @()
    if ($forceClean) {
        $buildArgs += "--no-cache", "--pull"
        Log-Info "🚫 캐시 없이 빌드합니다..."
    }
    
    # Backend 이미지 빌드
    Log-Info "Backend 이미지 빌드 중..."
    $buildCommand = @("build") + $buildArgs + @("-t", $BackendImage, "./iksoon_account_backend/")
    $result = & docker $buildCommand
    if ($LASTEXITCODE -ne 0) {
        Log-Error "Backend 이미지 빌드 실패"
        exit 1
    }
    Log-Success "Backend 이미지 빌드 완료"
    
    # Frontend 이미지 빌드
    Log-Info "Frontend 이미지 빌드 중..."
    $buildCommand = @("build") + $buildArgs + @("-t", $FrontendImage, "./iksoon_account_frontend/")
    $result = & docker $buildCommand
    if ($LASTEXITCODE -ne 0) {
        Log-Error "Frontend 이미지 빌드 실패"
        exit 1
    }
    Log-Success "Frontend 이미지 빌드 완료"
    
    Log-Success "모든 이미지 빌드 완료"
}

# 함수: 컨테이너 실행
function Start-Containers {
    Log-Info "컨테이너 시작 중..."
    docker-compose up -d
    
    if ($LASTEXITCODE -ne 0) {
        Log-Error "컨테이너 시작 실패"
        exit 1
    }
    
    # 헬스체크 대기
    Log-Info "서비스 시작 대기 중..."
    Start-Sleep -Seconds 10
    
    # Backend 헬스체크 (Frontend 프록시를 통해 확인)
    Log-Info "Backend 서비스 확인 중... (프록시를 통해)"
    $backendReady = $false
    for ($i = 1; $i -le 30; $i++) {
        try {
            # Frontend가 먼저 준비되어야 프록시 테스트 가능
            $frontendResponse = Invoke-WebRequest -Uri "http://localhost:3000" -TimeoutSec 2 -UseBasicParsing 2>$null
            if ($frontendResponse.StatusCode -eq 200) {
                # Frontend 프록시를 통해 Backend API 확인
                try {
                    $backendResponse = Invoke-WebRequest -Uri "http://localhost:3000/api/health" -TimeoutSec 2 -UseBasicParsing 2>$null
                    if ($backendResponse.StatusCode -eq 200) {
                        Log-Success "Backend 서비스가 정상적으로 시작되었습니다. (프록시 통신 확인)"
                        $backendReady = $true
                        break
                    }
                }
                catch {
                    Log-Warning "Frontend는 실행 중이지만 Backend API 프록시 연결 대기 중... ($i/30)"
                }
            }
            else {
                Log-Info "Frontend 서비스 대기 중... ($i/30)"
            }
        }
        catch {
            Log-Info "Frontend 서비스 대기 중... ($i/30)"
        }
        Start-Sleep -Seconds 2
    }
    
    if (-not $backendReady) {
        Log-Error "Backend 서비스 시작 실패 또는 프록시 연결 실패"
        Log-Info "디버깅 정보:"
        try {
            $frontendStatus = (Invoke-WebRequest -Uri "http://localhost:3000" -TimeoutSec 2 -UseBasicParsing).StatusCode
            Log-Info "  - Frontend 직접 접근: $frontendStatus"
        } catch {
            Log-Info "  - Frontend 직접 접근: FAIL"
        }
        try {
            $proxyStatus = (Invoke-WebRequest -Uri "http://localhost:3000/api/health" -TimeoutSec 2 -UseBasicParsing).StatusCode
            Log-Info "  - Backend 프록시 접근: $proxyStatus"
        } catch {
            Log-Info "  - Backend 프록시 접근: FAIL"
        }
        Log-Info "컨테이너 상태 확인:"
        docker-compose ps
        exit 1
    }
    
    # Frontend 헬스체크
    Log-Info "Frontend 서비스 확인 중..."
    $frontendReady = $false
    for ($i = 1; $i -le 30; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:3000" -TimeoutSec 2 -UseBasicParsing 2>$null
            if ($response.StatusCode -eq 200) {
                Log-Success "Frontend 서비스가 정상적으로 시작되었습니다."
                $frontendReady = $true
                break
            }
        }
        catch {
            # 아직 준비되지 않음
        }
        Start-Sleep -Seconds 2
    }
    
    if (-not $frontendReady) {
        Log-Error "Frontend 서비스 시작 실패"
        exit 1
    }
}

# 함수: 컨테이너 상태 확인
function Check-Status {
    Log-Info "컨테이너 상태 확인 중..."
    docker-compose ps
    
    Write-Output ""
    Log-Info "서비스 접속 정보:"
    Write-Output "  🌐 Frontend: http://localhost:3000"
    Write-Output "  🔧 Backend API: http://localhost:8080"
    Write-Output ""
    Log-Info "로그 확인: docker-compose logs -f"
    Log-Info "서비스 중지: docker-compose down"
}

# 함수: 사용법 출력
function Show-Usage {
    Write-Output "사용법: .\deploy.ps1 [옵션]"
    Write-Output ""
    Write-Output "옵션:"
    Write-Output "  -Clean        기존 이미지를 제거하고 새로 빌드"
    Write-Output "  -ForceClean   🧹 모든 Docker 캐시 및 빌드 캐시 완전 제거 후 빌드"
    Write-Output "  -Stop         실행 중인 컨테이너만 중지"
    Write-Output "  -Status       현재 실행 상태 확인"
    Write-Output "  -Help         이 도움말 출력"
    Write-Output ""
    Write-Output "예시:"
    Write-Output "  .\deploy.ps1                 # 일반 배포"
    Write-Output "  .\deploy.ps1 -Clean          # 이미지 제거 후 빌드"
    Write-Output "  .\deploy.ps1 -ForceClean     # 🧹 모든 캐시 제거 후 완전 새로 빌드"
    Write-Output "  .\deploy.ps1 -Stop           # 서비스 중지"
    Write-Output "  .\deploy.ps1 -Status         # 상태 확인"
    Write-Output ""
    Write-Output "⚠️  -ForceClean은 모든 Docker 캐시를 제거하므로 시간이 오래 걸릴 수 있습니다."
}

# 메인 실행 로직
function Main {
    Write-Output ""
    Log-Info "===== 스마트 가계부 Docker 배포 스크립트 ====="
    Write-Output ""
    
    if ($Help) {
        Show-Usage
        return
    }
    
    if ($Stop) {
        Cleanup-Containers
        Log-Success "서비스가 중지되었습니다."
        return
    }
    
    if ($Status) {
        Check-Status
        return
    }
    
    # 배포 옵션 확인
    Check-Docker
    
    if ($ForceClean) {
        # 완전 캐시 제거 배포
        Cleanup-Containers -cleanImages:$false -forceClean:$true
        Create-DataDirectory
        Build-Images -forceClean:$true
        Start-Containers
        Check-Status
    }
    elseif ($Clean) {
        # 일반 클린 배포
        Cleanup-Containers -cleanImages:$true
        Create-DataDirectory
        Build-Images
        Start-Containers
        Check-Status
    }
    else {
        # 일반 배포
        Cleanup-Containers
        Create-DataDirectory
        Build-Images
        Start-Containers
        Check-Status
    }
    
    Write-Output ""
    Log-Success "배포가 완료되었습니다! 🎉"
    Log-Info "브라우저에서 http://localhost:3000 에 접속하세요."
}

# 스크립트 실행
Main
