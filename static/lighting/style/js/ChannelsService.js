var module = angular.module('ChannelsService', ['LightingAPIService', 'LightingWebSocketService']);

module.factory('ChannelsService', ['$q', '$rootScope', 'LightingAPIService', 'LightingWebSocketService', function($q, $rootScope, LightingAPIService, LightingWebSocketService) {
	var channelsCachedPromise = null;

	var channelListeners = {};

	var currentChannelLevels = {};
	var lastSeenSeqNos = {};

	LightingWebSocketService.attachSocketListener('uC', function(data){
		var channelId = data.c.i;
		var universe = data.c.u;
		var channelLevel = data.c.l;
		var seqNo = data.c.s;

		$rootScope.$apply(function() {
			notifyChannelLister(universe, channelId, channelLevel, seqNo);
		});

	});

	function notifyChannelLister(universe, channelId, channelLevel, seqNo) {
		var lastSeenSeqNo = 0;
		var channelKey = universe + ":" + channelId;

		if (lastSeenSeqNos.hasOwnProperty(channelKey)) {
			lastSeenSeqNo = lastSeenSeqNos[channelKey];
		}
		if (seqNo > lastSeenSeqNo) {
			lastSeenSeqNos[channelKey] = seqNo;
			currentChannelLevels[channelKey] = channelLevel;

			if(channelListeners.hasOwnProperty(channelKey)) {
				var arr = channelListeners[channelKey];
				for (var i = 0; i < arr.length; i++) {
					arr[i](channelLevel);
				}
			}
		}
	}

	var Service = {};

	Service.getChannels = function(forceReload) {
		if (!channelsCachedPromise || forceReload) {
			channelsCachedPromise = LightingWebSocketService.channelState();

			channelsCachedPromise.then(function(channels) {
				currentChannelLevels = {};
				channels.forEach(function(channel) {
					var lastSeenSeqNo = 0;
					var channelKey = channel.universe + ":" + channel.id;

					if (lastSeenSeqNos.hasOwnProperty(channelKey)) {
						lastSeenSeqNo = lastSeenSeqNos[channelKey];
					}

					if (lastSeenSeqNo > channel.seqNo) {
						channel.currentLevel = currentChannelLevels[channelKey];
					} else {
						currentChannelLevels[channelKey] = channel.currentLevel;
						lastSeenSeqNos[channelKey] = channel.seqNo;
					}
				});
			});
		}

		return channelsCachedPromise;
	};

	Service.updateChannel = function(universe, id, level) {
		var channelKey = universe + ":" + id;
		if (lastSeenSeqNos.hasOwnProperty(channelKey)) {
			LightingWebSocketService.updateChannel(universe, id, level, lastSeenSeqNos[channelKey]);
		}
	};

	Service.attachChannelListener = function(universe, channelId, callback) {
		var channelKey = universe + ":" + channelId;
		if (!channelListeners.hasOwnProperty(channelKey)) {
			channelListeners[channelKey] = [];
		}
		channelListeners[channelKey].push(callback);
	};
	Service.detachChannelListener = function(universe, channelId, callback) {
		var channelKey = universe + ":" + channelId;
		if (channelListeners.hasOwnProperty(channelKey)) {
			var array = channelListeners[channelKey];
			var index = array.indexOf(callback);

			if (index > -1) {
				array.splice(index, 1);
			}
		}
	};

	Service.currentChannelLevel = function(universe, channelId) {
		var channelKey = universe + ":" + channelId;
		return currentChannelLevels[channelKey];
	}

	return Service;
}]);
