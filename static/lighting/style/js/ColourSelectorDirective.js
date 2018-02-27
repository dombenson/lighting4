var module = angular.module('ColourSelectorDirective', ['ChannelsService']);

module.directive('colourSelector', function() {
	return {
		restrict: 'A',
		templateUrl: '/lighting/style/partials/colourSelector.html',
		scope: {
			redChannel:   "=redchannel",
			greenChannel: "=greenchannel",
			blueChannel:  "=bluechannel",
			universe:     "=universe"
		},
		controller: function($scope, ChannelsService) {
			$scope.currentColour = "#ffffff";

			var currentRedLevel   = ChannelsService.currentChannelLevel($scope.universe, $scope.redChannel)     || 0;
			var currentGreenLevel = ChannelsService.currentChannelLevel($scope.universe, $scope.greenChannel) || 0;
			var currentBlueLevel  = ChannelsService.currentChannelLevel($scope.universe, $scope.blueChannel)    || 0;

			updateCurrentMode();

			var redListener = function(newLevel){
				currentRedLevel = newLevel;
				updateCurrentMode();
			};
			var greenListener = function(newLevel){
				currentGreenLevel = newLevel;
				updateCurrentMode();
			};
			var blueListener = function(newLevel){
				currentBlueLevel = newLevel;
				updateCurrentMode();
			};

			ChannelsService.attachChannelListener($scope.universe, $scope.redChannel, redListener);
			ChannelsService.attachChannelListener($scope.universe, $scope.greenChannel, greenListener);
			ChannelsService.attachChannelListener($scope.universe, $scope.blueChannel, blueListener);

			$scope.$on('$destroy', function() {
				ChannelsService.detachChannelListener($scope.universe, $scope.redChannel, redListener);
				ChannelsService.detachChannelListener($scope.universe, $scope.greenChannel, greenListener);
				ChannelsService.detachChannelListener($scope.universe, $scope.blueChannel, blueListener);
			});

			function updateCurrentMode() {
				$scope.currentColour = rgbToHex(currentRedLevel, currentGreenLevel, currentBlueLevel);
			}

			function componentToHex(c) {
				var hex = c.toString(16);
				return hex.length == 1 ? "0" + hex : hex;
			}

			function rgbToHex(r, g, b) {
				return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
			}
		}
	};
});