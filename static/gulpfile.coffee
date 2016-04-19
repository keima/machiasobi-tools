gulp = require 'gulp'
$ = require('gulp-load-plugins')() # injecting gulp-* plugin
browserSync = require 'browser-sync'
bowerFiles = require "main-bower-files"
runSequence = require 'run-sequence'
rimraf = require "rimraf"

# config
config =
  js: 'scripts/**/*.js'
  css: 'stylesheets/**/*.css'
  partial: 'partials/**/*.html'

  output: '../src/static/'

  watching: false

#
# Task
#
gulp.task 'browser-sync', ->
  browserSync proxy: 'localhost:8080'

gulp.task "setWatch", ->
  config.watching = true

gulp.task 'watch', ["setWatch"], ->
  gulp.watch [config.js, config.css], ['inject', 'usemin', browserSync.reload]
  gulp.watch config.partial, ['copyPartials']

gulp.task 'inject', ->
  bower = gulp.src bowerFiles(), {read: false}
  js = gulp.src([config.js]).pipe($.angularFilesort())
  css = gulp.src [config.css], {read: false}

  gulp.src './index.html'
  .pipe $.inject bower, {addRootSlash: false, name: 'bower'}
  .pipe $.inject js, {addRootSlash: false}
  .pipe $.inject css, {addRootSlash: false}
  .pipe gulp.dest './'

gulp.task 'usemin', ->
  cssTask = (files, filename) ->
    s = files
    s = s.pipe $.cleanCss() if !config.watching
    s = s.pipe $.concat(filename)
    s = s.pipe $.rev() if !config.watching
    return s

  jsTask = (files, filename) ->
    s = files
    s = s.pipe $.ngAnnotate() if !config.watching
    s = s.pipe $.uglify() if !config.watching
    s = s.pipe $.concat(filename)
    s = s.pipe $.rev() if !config.watching
    return s

  gulp.src './index.html'
  .pipe $.spa.html(
    pipelines:
      main: (files)->
        s = files
        s = s.pipe $.minifyHtml(empty: true, conditionals: true) if !config.watching
        return s
      css: (files)->
        cssTask files, "app.css"
      vendorjs: (files)->
        jsTask files, "vendor.js"
      js: (files)->
        jsTask files, "app.js"
  )
  .pipe gulp.dest(config.output)

gulp.task 'copy', ['copyPartials', 'copyOthers']

gulp.task 'copyPartials', ->
  gulp.src config.partial, {base: './'}
  .pipe $.minifyHtml(empty: true)
  .pipe gulp.dest config.output

gulp.task 'copyOthers', ->
  # other
  gulp.src ['*.png'], {base: './'}
  .pipe gulp.dest config.output


gulp.task 'clean', (cb) ->
  rimraf(config.output, cb);


gulp.task 'default', ['build', 'browser-sync', 'watch']

gulp.task 'build', (cb) -> runSequence(
  'clean', 'inject', ['usemin', 'copy']
  cb
)
