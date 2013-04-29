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
	'localhost:3000',
	'localhost:3001'
];

consumers = [
	{id: 'Pepe'},
	{id: 'Juan'}
];

services = [
	{id: 'Service1'},
	{id: 'Service2'},
	{id: 'Service3'},
	{id: 'Service4'}
];

servers = [
	{name: 'server1.com', services: ['Service1']},
	{name: 'server2.com', services: ['Service3']},
	{name: 'server3.com', services: ['Service1','Service2','Service3']},
];

active_links = []

module.exports = {
	err_enum: err_enum,
	err_back: err_back,
	siblings: siblings,
	consumers: consumers,
	services: services,
	servers: servers,
	active_links: active_links,
};