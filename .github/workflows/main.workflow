workflow "Deploy to Heroku" {
  on = "push"
  resolves = ["release registry", "release hooks", "release reaper"]
}

action "only master branch" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "heroku login" {
  needs = "only master branch"
  uses = "actions/heroku@master"
  args = "container:login"
  secrets = ["HEROKU_API_KEY"]
}

action "build registry" {
  uses = "actions/docker/cli@master"
  needs = "heroku login"
  args = "build -t registry.heroku.com/ttlsh/web registry"
}

action "push registry" {
  uses = "actions/docker/cli@master"
  needs = "build registry"
  args = "push registry.heroku.com/ttlsh/web"
}

action "release registry" {
  uses = "actions/heroku@master"
  needs = "push registry"
  args = "container:release -a ttlsh web"
  secrets = ["HEROKU_API_KEY"]
}

action "build hooks" {
  uses = "actions/docker/cli@master"
  needs = "heroku login"
  args = "build -f hooks/Dockerfile.hooks -t registry.heroku.com/ttlsh-hooks/web hooks"
}

action "push hooks" {
  uses = "actions/docker/cli@master"
  needs = "build hooks"
  args = "push registry.heroku.com/ttlsh-hooks/web"
}

action "release hooks" {
  uses = "actions/heroku@master"
  needs = "push hooks"
  args = "container:release -a ttlsh-hooks web"
  secrets = ["HEROKU_API_KEY"]
}

action "build reaper" {
  uses = "actions/docker/cli@master"
  needs = "heroku login"
  args = "build -f hooks/Dockerfile.reap -t registry.heroku.com/ttlsh-hooks/reap hooks"
}

action "push reaper" {
  uses = "actions/docker/cli@master"
  needs = "build reaper"
  args = "push registry.heroku.com/ttlsh-hooks/reap"
}

action "release reaper" {
  uses = "actions/heroku@master"
  needs = "push reaper"
  args = "container:release -a ttlsh-hooks reap"
  secrets = ["HEROKU_API_KEY"]
}
