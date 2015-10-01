gulp = require 'gulp'
$ = require('gulp-load-plugins')() # injecting gulp-* plugin
browserSync = require 'browser-sync'
bowerFiles = require "main-bower-files"
runSequence = require 'run-sequence'
rimraf = require "rimraf"

# config
config =
  js: './scripts/**/*.js'
  css: './stylesheets/**/*.css'
  output: '../static-build/'

#
# Task
#
gulp.task 'browser-sync', ->
  browserSync proxy: 'localhost:8080'

gulp.task 'watch', ->
  gulp.watch config.js, ['inject', browserSync.reload]
  gulp.watch config.css, ['inject', browserSync.reload]

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
  gulp.src './index.html'
  .pipe $.usemin(
    css: [$.minifyCss(), $.rev()]
    html: [$.minifyHtml(empty: true, conditionals: true)]
    js: [$.ngAnnotate(), $.uglify(), $.rev()]
    vendorjs: [$.ngAnnotate(), $.uglify(), $.rev()]
  )
  .pipe gulp.dest(config.output)

gulp.task 'copy', ->
  gulp.src './partials/**/*.html', {base: './'}
  .pipe $.minifyHtml(empty: true)
  .pipe gulp.dest config.output

  # other
  gulp.src ['*.png'], {base: './'}
  .pipe gulp.dest config.output


gulp.task 'clean', (cb) ->
  rimraf(config.output, cb);


gulp.task 'default', ['browser-sync', 'watch']

gulp.task 'build', (cb) -> runSequence(
  'clean', 'inject', 'usemin', 'copy'
  cb
)
