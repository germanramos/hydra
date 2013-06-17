var enums = module.exports;

enums.app = {
	localStrategyEnum : {
		INDIFFERENT: 0,
		ROUND_ROBIN: 1,
		SERVER_LOAD: 2,
		CHEAPEST: 3
	},

	cloudStrategyEnum : {
		INDIFFERENT: 0,
		ROUND_ROBIN: 1,
		CHEAPEST: 2,
		CLOUD_LOAD: 3
	},

	stateEnum : {
		READY: 0,
		UNAVAILABLE : 1,
		DELETE: 2
	}
};

enums.server = {
	stateEnum : {
		READY: 0,
		UNAVAILABLE : 1
	}
};