err_enum = {
    ok : 'ok',
    not_here : 'not_here',
    wrong_service : 'wrong_service',
    wrong_consumer : 'wrong_consumer'
};

err_back = {
    ok : 'ok',
    already_started : 'already_started',
    already_stopped : 'already_stopped',
    denied : 'denied'
};

siblings = [
	'http://localhost:3000',
	//'http://localhost:3001'
];

consumers = [
	{id: 'Pepe'},
	{id: 'Juan'}
];

services = [
	{id: 'service_time'},
	{id: 'service_sum'},
	{id: 'service_uuid'},
	{id: 'service_other'}
];

servers = [
	{name: 'http://localhost:4001', services: ['service_time']},
	{name: 'http://localhost:4003', services: ['service_uuid']},
	{name: 'http://localhost:4002', services: ['service_sum']}
];

active_links = [];

module.exports = {
	err_enum: err_enum,
	err_back: err_back,
	siblings: siblings,
	consumers: consumers,
	services: services,
	servers: servers,
	active_links: active_links
};