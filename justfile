#!/usr/bin/env just --justfile
# list of available commands
default:
  @just --list

regFile := "./registries"
tagCommit := `[ -d ".git" ] && (git rev-list --abbrev-commit --tags --max-count=1) || true`
lastVersion := `[ -d ".git" ] && (git describe --tags --abbrev=0 2>/dev/null || git symbolic-ref -q --short HEAD) || true`
lastBranchName := `[ -d ".git" ] && (git describe --tags --exact-match HEAD 2>/dev/null || git symbolic-ref -q --short HEAD) || true`
goVersion := `if [ -f "go.mod" ]; then grep -oP '^go\s+\K\d+(\.\d+)?' go.mod; else go version | sed -n 's/.*go\([0-9.]*\).*/\1/p'; fi`
serviceName := `if [ -f go.mod ]; then grep '^module' go.mod | awk '{print $2}' | sed -E 's#.*/([^/]+)\.git#\1#; s#.*/([^/]+)#\1#'; elif [ -d .git ]; then git remote get-url origin 2>/dev/null | sed -E 's#.*/([^/]+)\.git#\1#; s#.*/([^/]+)#\1#'; fi`
repoRemoteURL := `[ -d ".git" ] && (git config --get remote.origin.url) || true`
githubOrigin := `[ -d ".git" ] && (git config --get remote.origin.url | grep github.com || echo "") || true`

# prints package informations
info:
    @printf "\033[1mPackage information:\033[0m\n"
    @printf "%-12s\t%s\n" "url:" "{{repoRemoteURL}}"
    @printf "%-12s\t%s\n" "service:" "{{serviceName}}"
    @printf "%-12s\t%s\n" "go-version:" "{{goVersion}}"
    @printf "%-12s\t%s\n" "current-tag:" "{{lastBranchName}}"


# clean build directory
clean:
    @echo "clean bin directory..."
    @[ -d "./bin" ] && rm -r ./bin && echo "bin directory cleaned" || true

# clean and build project
build: clean
    go build -o ./bin/service -ldflags="-s -w" -ldflags="-X 'main.Version={{lastBranchName}}' -X 'main.BuildDate=$(date -u '+%Y-%m-%d %H:%M:%S')'" ./cmd

# build and compress binary
upx: build
    upx --best --lzma bin/service


# build service docker image
buildimg tag="auto" cache="true" latest="true":
    #!/usr/bin/env bash
    echo "build service docker image..."
    just info
    imgtag="{{tag}}"
    if [[ "$imgtag" == "auto" ]];then
        imgtag="{{lastBranchName}}"
    fi
    imgname="ghcr.io/traxtex/{{serviceName}}"
    echo "building image: $imgname:$imgtag"
    if [[ "{{cache}}" == "true" ]];then
        docker buildx build -t "$imgname:$imgtag" -f Dockerfile --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" --build-arg GO_VERSION="{{goVersion}}" --build-arg GITHUB_ORIGIN="{{githubOrigin}}" .
    else
        docker buildx build -t "$imgname:$imgtag" -f Dockerfile --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" --build-arg GO_VERSION="{{goVersion}}" --build-arg GITHUB_ORIGIN="{{githubOrigin}}" --no-cache .
    fi
    if [[ "{{ latest }}" == "true" ]];then
        echo "building image: $imgname:latest"
        docker tag "$imgname:$imgtag" "$imgname:latest"
    fi