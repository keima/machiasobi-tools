gulp = require 'gulp'
p = require('gulp-load-plugins')() # injecting gulp-* plugin
browserSync = require 'browser-sync'
reload = browserSync.reload

# config
config =
  js: './scripts/**/*.js'
  css: './stylesheets/**/*.css'

gulp.task 'inject', ->
  sources = gulp.src [config.js, config.css], {read: false}

  gulp.src './index.html'
  .pipe p.inject sources
  .pipe gulp.dest './'

gulp.task 'browser-sync', ->
  browserSync proxy: 'localhost:8080'

gulp.task 'default', ['browser-sync'], ->
  gulp.watch config.js, ['inject', reload]
  gulp.watch config.css, ['inject', reload]