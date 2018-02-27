var module = angular.module('ChannelGroupDirective', ['ChannelsService']);

module.directive('channelGroup', function(ChannelsService) {
	return {
		restrict: 'A',
		templateUrl: '/lighting/style/partials/channelGroup.html',
		scope: {
			channelGroup: "=channelgroup",
			firstChannel: "=firstchannel",
			universe: "=universe"
		},
		controller: function($scope) {
			$scope.toggleDisabled = function(channel) {
				if (channel.disabled) {
					ChannelsService.updateChannel($scope.universe, channel.channelOffset + $scope.firstChannel, channel.disabledDefaultLevel);
				} else {
					ChannelsService.updateChannel($scope.universe, channel.channelOffset + $scope.firstChannel, channel.enabledDefaultLevel);
				}
			};
		}
	};
});