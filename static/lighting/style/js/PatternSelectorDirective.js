var module = angular.module('PatternSelectorDirective', ['ChannelsService']);

module.directive('patternSelector', function() {
	return {
		restrict: 'A',
		templateUrl: '/lighting/style/partials/patternSelector.html',
		scope: {
			groupChannel:   "=groupchannel",
			patternChannel: "=patternchannel",
			patterns:       "=patterns",
			universe:       "=universe"
		},
		controller: function($scope, ChannelsService) {
			var currentGroupLevel = ChannelsService.currentChannelLevel($scope.universe, $scope.groupChannel) || 0;
			var currentPatternLevel = ChannelsService.currentChannelLevel($scope.universe, $scope.patternChannel) || 0;

			$scope.currentPattern = 'NONE';

			$scope.selectPattern = function(pattern) {
				ChannelsService.updateChannel($scope.universe, $scope.groupChannel, pattern.groupLevel);
				ChannelsService.updateChannel($scope.universe, $scope.patternChannel, pattern.patternLevel);
			}

			calculateCurrentPattern();

			var groupListener = function(newLevel){
				currentGroupLevel = newLevel;
				calculateCurrentPattern();
			};
			var patternListener = function(newLevel){
				currentPatternLevel = newLevel;
				calculateCurrentPattern();
			};
			ChannelsService.attachChannelListener($scope.universe, $scope.groupChannel, groupListener);
			ChannelsService.attachChannelListener($scope.universe, $scope.patternChannel, patternListener);

			$scope.$on('$destroy', function() {
				ChannelsService.detachChannelListener($scope.universe, $scope.groupChannel, groupListener);
				ChannelsService.detachChannelListener($scope.universe, $scope.patternChannel, patternListener);
			});

			function calculateCurrentPattern() {
				var currentSelection = null;
				$scope.patterns.forEach(function(pattern) {
					if (pattern.groupLevel <= currentGroupLevel && pattern.patternLevel <= currentPatternLevel) {
						if (currentSelection == null ||
										(pattern.groupLevel >= currentSelection.groupLevel && pattern.patternLevel >= currentSelection.patternLevel)) {
							currentSelection = pattern;
						}
					}
				});
				if (currentSelection == null) {
					$scope.currentPattern = 'NONE';
				} else {
					$scope.currentPattern = currentSelection.name;
				}
			}
		}
	};
});