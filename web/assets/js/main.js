/* globals jQuery, Backbone */
(function($) {
	'use strict';

	var
		apiVersion = 'v1',

		// Model and Collection.
		ChallengeModel,
		CurrentChallengeModel,
		ChallengeCollection,
		UserModel,
		CurrentUserModel,
		SubmissionModel,
		SubmissionCollection,

		// Views.
		CurrentChallengeView,
		PastChallengesView,
		PastChallengeItem,
		SubmissionsView,
		SubmissionItem,

		// Instances.
		CurrentChallenge,
		Challenges,

		// Router.
		Router,

		// URL for Model and Collection.
		getAPIURL = function(suffix) {
			return '/' + apiVersion + suffix;
		};

	// ChallengeModel
	// --------------
	//
	// ChallengeModel represents a single challenge.
	ChallengeModel = Backbone.Model.extend({
		urlRoot: function() {
			return getAPIURL('/challenges');
		},
		parse: function(model) {
			// `import` is reserved word.
			if (model.import) {
				model.import_url = model.import;
			}
			return model;
		}
	});

	// CurrentChallengeModel
	// ---------------------
	//
	// CurrentChallengeModel represents current challenge.
	CurrentChallengeModel = Backbone.Model.extend({
		url: getAPIURL('/challenges/current')
	});

	// CurrentChallenge
	// ----------------
	//
	// Instance of CurrentChallengeModel.
	CurrentChallenge = new CurrentChallengeModel();

	// ChallengeCollection
	// -------------------
	//
	// ChallengeCollection represents challenges collection.
	ChallengeCollection = Backbone.Collection.extend({
		model: ChallengeModel,
		url: function() {
			return getAPIURL('/challenges');
		}
	});

	// Challenges
	// ----------
	//
	// Instance of ChallengeCollection.
	Challenges = new ChallengeCollection();

	// UserModel
	// ---------
	//
	// UserModel represents user. Can be participant or evaluator.
	UserModel = Backbone.Model.extend({
		defaults: {
			name: '',
			avatar_url: '',
			email: ''
		},
		urlRoot: function() {
			return getAPIURL('/users');
		}
	});

	// CurrentUserModel
	// ----------------
	//
	// CurrentUserModel represents current user.
	CurrentUserModel = UserModel.extend({
		url: function() {
			return getAPIURL('/user');
		}
	});

	// SubmissionModel
	// ---------------
	//
	// SubmissionModel represents challenge submission.
	SubmissionModel = Backbone.Model.extend();

	// SubmissionCollection
	// --------------------
	//
	// SubmissionCollection represents all submissions of a challenge.
	SubmissionCollection = Backbone.Collection.extend({
		model: SubmissionModel,
		challengeID: null,

		initialize: function(options) {
			this.challengeID = options.challengeID;
		},

		url: function() {
			return getAPIURL('/challenges/' + this.challengeID + '/submissions');
		}
	});

	// CurrentChallengeView
	// --------------------
	//
	// DOM Element for current challenge in home page.
	CurrentChallengeView = Backbone.View.extend({
		initialize: function() {
			this.template = _.template($('#current-challenge-tpl').html());
			this.content = this.$('.current-challenge-content');

			this.listenTo(CurrentChallenge, 'sync', this.render);
		},

		render: function(model) {
			this.content.html(this.template(model.toJSON()));
			return this;
		}
	});

	// PastChallengesView
	// ------------------
	//
	// DOM Element for past challenges in home page.
	PastChallengesView = Backbone.View.extend({
		initialize: function() {
			this.listenTo(Challenges, 'reset', this.render);
		},

		render: function() {
			Challenges.each(this.renderItem, this);
		},

		renderItem: function(model) {
			if (model.get('status') === 'open') return;

			var view = new PastChallengeItem({model: model});
			this.$el.append(view.render().el);
		}
	});

	// PastChallengeItem
	// -----------------
	//
	// DOM Element of a single challenge in PastChallengesView
	PastChallengeItem = Backbone.View.extend({
		tagName: 'li',
		className: 'challenge',

		initialize: function() {
			this.template = _.template($('#past-challenge-item-tpl').html());
		},

		render: function() {
			this.$el.html(this.template(this.model.toJSON()));
			return this;
		}
	});

	// SubmissionsView
	// ---------------
	//
	// DOM Element for submissions list.
	SubmissionsView = Backbone.View.extend({
		render: function() {
			this.collection.each(this.renderItem, this);
		},

		renderItem: function(model) {
			var view = new SubmissionItem({model: model});
			this.$el.append(view.render().el);
		}
	});

	// SubmissionItem
	// --------------
	//
	// DOM Element for single submission in submissions list.
	SubmissionItem = Backbone.View.extend({
		tagName: 'li',
		className: 'submission',

		initialize: function() {
			this.template = _.template($('#submission-item-tpl').html());
		},

		render: function() {
			console.log(this.model.toJSON());
			this.$el.html(this.template(this.model.toJSON()));
			return this;
		}
	});

	// Router
	// ------
	Router = Backbone.Router.extend({
		routes: {
			"challenges/:challenge_id":  "viewChallenge",
			"api_key=*fragment":         "getAPIKeyFromURL",
			"":                          "home",
			"profile":                   "viewProfile",
			"logout":                    "logout",
			"*path":                     "home"
		},

		initialize: function() {
			this.userNavTemplate = _.template($('#user-nav-tpl').html());
			this.profileTemplate = _.template($('#profile-tpl').html());
			this.homeTemplate = _.template($('#home-tpl').html());
			this.challengeTemplate = _.template($('#challenge-tpl').html());

			this.userNav = $('.user-nav');
			this.content = $('#content');

			this.currentUser = null;
			this.renderUserNav();
		},

		home: function() {
			this.content.html(this.homeTemplate());
			this.currentChallengeView = new CurrentChallengeView({
				el: $('#current-challenge')
			});
			this.pastChallengesView = new PastChallengesView({
				el: $('#past-challenges')
			});

			Challenges.fetch({reset: true});
			CurrentChallenge.fetch();
		},

		viewChallenge: function(challenge_id) {
			var challenge = new ChallengeModel({
				id: challenge_id
			});

			challenge.once('sync', this._challengeRender, this);
			challenge.fetch();
		},

		_challengeRender: function(challenge) {
			this.content.html(this.challengeTemplate(challenge.toJSON()));
			this._challengeSubmissions(challenge);
		},

		_challengeSubmissions: function(challenge) {
			var submissions = new SubmissionCollection({
				challengeID: challenge.get('id')
			});

			submissions.once('sync', this._challengeSubmissionsRender, this);
			submissions.fetch();
		},

		_challengeSubmissionsRender: function(submissions) {
			if (submissions.size() === 0) {
				$('.challenge-submissions').html(
					'<li><p>No submission for this challenge</p></li>'
				);
			} else {
				new SubmissionsView({
					collection: submissions,
					el: $('.challenge-submissions')
				}).render();
			}
		},

		getAPIKeyFromURL: function() {
			var parts = location.hash.slice(1).split('='),
				api_key = '';

			if (parts[0] === 'api_key') {
				api_key = parts[1];
			}
			if (api_key !== '') {
				localStorage.setItem("api_key", decodeURIComponent(api_key));
				this.renderUserNav();
			}

			this.navigate("", {trigger: true});
		},

		renderUserNav: function() {
			var key = this.getAPIKey();
			if (!key) {
				this.userNav.html(this.userNavTemplate({user: null}));
				return false;
			}

			this.currentUser = new CurrentUserModel();
			this.currentUser.once('sync', function(){
				this.userNav.html(this.userNavTemplate({
					user: this.currentUser.toJSON()
				}));
			}, this);

			this.fetchCurrentUser();
		},

		getAPIKey: function() {
			return localStorage.getItem("api_key");
		},

		fetchCurrentUser: function() {
			this.currentUser.fetch({
				beforeSend: $.proxy(this, 'setHeader'),
				error: $.proxy(this, 'fetchCurrentUserError')
			});
		},

		setHeader: function(xhr) {
			xhr.setRequestHeader("Auth-ApiKey", this.getAPIKey());
		},

		fetchCurrentUserError: function() {
			this.logout();
		},

		viewProfile: function() {
			if (!this.currentUser) {
				this.navigate("", {trigger: true});
				return false;
			}

			this.currentUser.once('sync', function() {
				var user = this.currentUser.toJSON();
				user.api_key = this.getAPIKey();
				this.content.html(this.profileTemplate({
					user: user
				}));
			}, this);
			this.fetchCurrentUser();
		},

		logout: function() {
			this.currentUser = null;
			localStorage.removeItem('api_key');

			this.renderUserNav();
			this.navigate("", {trigger: true});
		}
	});

	// Start router when DOM is ready.
	$(function() {
		new Router();
		Backbone.history.start();
	});

}(jQuery));
