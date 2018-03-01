var lightingApp = angular.module('lightingApp', ['ChannelSliderDirective', 'PatternSelectorDirective', 'ChannelGroupDirective', 'ColourSelectorDirective', 'ModeSelectorDirective', 'ModeToggleDirective', 'LightingWebSocketService', 'ChannelsService', 'FixtureService', 'FixtureDirective', 'SequenceService', 'SequenceDirective', 'SequenceStepDirective', 'TrackService', 'PlaybackDetailsDirective']);

lightingApp.controller('WebSocketManagerController', function($scope, LightingWebSocketService, ChannelsService, FixtureService, SequenceService) {
	$scope.currentView = 'channels';
	$scope.currentUniverse = 0;

	$scope.currentWebSocketStatus = 'striped';
	$scope.canReconnect = false;
	$scope.connected = false;

	$scope.universes = {};
	$scope.sequences = [];

	$scope.reconnect = function(){
		LightingWebSocketService.reconnect();
		$scope.canReconnect = false;
	};

	$scope.blackout = function() {
		LightingWebSocketService.sendRequest('blackout');
	};
	$scope.unblackout = function() {
		LightingWebSocketService.sendRequest('unblackout');
	};

	$scope.sequenceBeat = function() {
		SequenceService.sequenceBeat();
	};

	LightingWebSocketService.attachConnectionListener(function(status) {
		if (status == 'open') {
			$scope.currentWebSocketStatus = 'successful';
			$scope.canReconnect = false;
			$scope.connected = true;
			loadInitialState();
		} else if (status == 'close' || status == 'error') {
			$scope.currentWebSocketStatus = 'failed';
			$scope.canReconnect = true;
			$scope.connected = false;
			$scope.universes = {};
		} else {
			$scope.currentWebSocketStatus = 'striped';
			$scope.canReconnect = false;
			$scope.connected = false;
			$scope.universes = {};
		}

	});

	function loadInitialState() {
		ChannelsService.getChannels(true).then(
			function(channels) {
				$scope.universes = channelsIntoUniverses(channels);

				if (!$scope.currentUniverse && $scope.universes) {
					$scope.currentUniverse = $scope.universes[Object.keys($scope.universes)[0]].number;
				}

				FixtureService.getFixtureList(true).then(
					function(fixtures) {
						fixturesIntoUniverses(fixtures, $scope.universes);
					}, function(error) {
						console.error('Could not load fixtures', error);
					}
				);
			}, function(error) {
				console.error('Could not load channels', error);
			}
		);
	}

	function fixturesIntoUniverses(fixtures, universes) {
		for (var universeKey in fixtures) {
			universes[universeKey].fixtures = fixtures[universeKey];
		}
	}

	function channelsIntoUniverses(channels) {
		var universes = {};

		for (var i = 0; i < channels.length; i++) {
			var channel = channels[i];

			if (!universes.hasOwnProperty(channel.universe)) {
				universes[channel.universe] = {
					number: channel.universe,
					channels: []
				};
			}

			universes[channel.universe].channels.push(channel);
		}

		for (var universeKey in universes) {
			var universe = universes[universeKey];
			universe.channelsInPages = channelsIntoPages(universe.channels, 8);
		}

		return universes;
	}

	function channelsIntoPages(channels, channelsPerPage) {
		var fullPageCount = Math.floor(channels.length / channelsPerPage);
		var finalPageChannels = channels.length % channelsPerPage;

		var pages = [];

		for (var i = 0; i < fullPageCount; i++) {
			var page = [];

			for (var j = 0; j < channelsPerPage; j++) {
				var channel = channels[i * channelsPerPage + j];
				page.push(channel);
			}

			pages.push(page);
		}

		if (finalPageChannels > 0) {
			var finalPage = [];
			for (var i = 0; i < finalPageChannels; i++) {
				var channel = channels[fullPageCount * channelsPerPage + i];
				finalPage.push(channel);
			}
			pages.push(finalPage);
		}

		return pages;
	}
});
