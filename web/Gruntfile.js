module.exports = function(grunt) {
	'use strict';

	grunt.initConfig({
		clean: ['dist/'],

		jshint: ['assets/js/main.js'],

		uglify: {
			dist: {
				files: {
					'assets/js/all.min.js': [
						'assets/js/libs/jquery/dist/jquery.min.js',
						'assets/js/libs/bootstrap/dist/js/bootstrap.min.js',
						'assets/js/libs/underscore/underscore-min.js',
						'assets/js/libs/backbone/backbone.js',
						'assets/js/main.js'
					]
				}
			}
		},

		cssmin: {
			dist: {
					files: {
					'assets/css/all.min.css': [
						'assets/js/libs/bootstrap/dist/css/bootstrap.min.css',
						'assets/css/main.css'
					]
				}
			}
		},

		processhtml: {
			dist: {
				files: {
					'dist/index.html': ['layout.html']
				}
			}
		},

		copy: {
			dist: {
				files: [
					{src: ['assets/img/**'], dest: 'dist/'}
				]
			}
		},

		watch: {
			js: {
				files: ['assets/js/*.js', '!assets/js/*.min.js', '!assets/js/libs/*'],
				tasks: ['default']
			},
			css: {
				files: ['assets/css/*.css', '!assets/css/*.min.css'],
				tasks: ['default']
			},
			img: {
				files: ['assets/img/*'],
				tasks: ['default']
			},
			html: {
				files: ['layout.html', 'templates/*.html'],
				tasks: ['default']
			}
		}
	});

	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-contrib-cssmin');
	grunt.loadNpmTasks('grunt-processhtml');
	grunt.loadNpmTasks('grunt-contrib-copy');
	grunt.loadNpmTasks('grunt-contrib-watch');

	grunt.registerTask('default', [
		'clean',
		'jshint',
		'uglify',
		'cssmin',
		'processhtml',
		'copy'
	]);
};
