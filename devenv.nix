{pkgs, ...}: {
  # https://devenv.sh/basics/

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    nodejs_22
    uv
    delve
    alejandra
    nixd
    just
    nodejs_22
    uv
    go
    gopls
    gotools
    go-tools
  ];

  env = {
    TOKEN_FILE = "./testtoken.txt";
    CONTAINERS = "[\"docker.io/library/alpine:latest\"]";
    CONTAINER_NAMES = "[\"container-refresh-test\"]";
    PORT = "9080";
    SLACK_WEBHOOK_URL_FILE = "./testwebhookurl.txt"; # Or comment out to disable
  };

  # https://devenv.sh/languages/
  # languages.rust.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.refresh-test.exec = ''
    curl -X POST -H "Content-Type: application/json" -d '{"token": "testtoken"}' http://localhost:9080/update
  '';

  scripts.refresh-test-run.exec = ''
    echo "Starting persistent container named container-refresh-test"
    podman run -d --name container-refresh-test alpine top
    echo "Started.  If you have the app running, run 'refresh-test' which should stop the container."
  '';

  enterShell = ''
    go version
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
