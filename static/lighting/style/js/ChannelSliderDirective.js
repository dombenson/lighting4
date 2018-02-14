var module = angular.module('ChannelSliderDirective', ['ChannelsService', 'TrackService']);

module.directive('channelSlider', function () {
	return {
		restrict: 'E',
		templateUrl: '/lighting/style/partials/channelSlider.html',
		scope: {
			channelid: '=channelid',
			maxlevel: '=?maxlevel',
			minlevel: '=?minlevel',
			enabled: '=?enabled',
			reverse: '=?reverse'
		},
		controller: function ($scope, ChannelsService, TrackService) {
			$scope.currentLevel = ChannelsService.currentChannelLevel($scope.channelid) || 0;
			$scope.minlevel = $scope.minlevel || 0;
			$scope.maxlevel = $scope.maxlevel || 255;

			if ($scope.enabled === undefined) {
				$scope.enabled = true;
			}

			if ($scope.reverse) {
				$scope.currentLevel = $scope.maxlevel - $scope.currentLevel;
			}

			var channelListener = function (newLevel) {
				if ($scope.reverse) {
					$scope.currentLevel = $scope.maxlevel - newLevel;
				} else {
					$scope.currentLevel = newLevel;
				}
			};

			ChannelsService.attachChannelListener($scope.channelid, channelListener);

			$scope.$watch("currentLevel", function (n, o) {
				if (n != o) {
					if ($scope.reverse) {
						ChannelsService.updateChannel($scope.channelid, $scope.maxlevel - parseInt($scope.currentLevel));
					} else {
						ChannelsService.updateChannel($scope.channelid, parseInt($scope.currentLevel));
					}
				}
			});

			$scope.$on('$destroy', function () {
				ChannelsService.detachChannelListener($scope.channelid, channelListener);
			});
		}
	};
});