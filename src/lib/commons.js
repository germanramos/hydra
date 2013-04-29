module.exports.request	= require('request');    	// nodejs request
module.exports.express	= require('express');    	// basic framework
module.exports.async 	= require('async');    		// asynchronous management
module.exports.q 		= require('q');      		// defered & promises management
module.exports.rabbit	= require("rabbit.js");  	// mq
module.exports.amqp		= require("amqp");  		// raw mq
module.exports.mongodb	= require("mongodb");    	// document data store
module.exports.redis	= require("redis");      	// key-value data store
module.exports.RedisStore = require('connect-redis')(module.exports.express); //redis session storage
module.exports.assert	= require("assert");     	// unit test
module.exports.fs		= require('fs');          	// file system
module.exports.extend	= require('xtend');			// merge object properties
module.exports.passport	= require('passport');		// security strategies framework
module.exports.GoogleStrategy  = require('passport-google-oauth').OAuthStrategy;
module.exports.FacebookStrategy= require('passport-facebook').Strategy;
module.exports.TwitterStrategy = require('passport-twitter').Strategy;
//module.exports.hero		= require('./hero.js');		// hero library
module.exports.hero		= require('hero');		// hero library
module.exports.request		= require('request');		// hero library
module.exports.http		= require('http');
module.exports.uuid = require('node-uuid');
module.exports.utils = require('connect').utils;
module.exports.pastry = require('pastry'); //Sokobank cookie manager
module.exports.qlog_node = require('qlog-node'); //Node module for Qlog