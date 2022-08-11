// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package extension

import (
	"aws-lambda-extensions/cache-extension-demo/plugins"
	"os"
	"strconv"
	"strings"
)

// Constants definition
const (
	Parameters = "parameters"
	InitializeCacheOnStartup = "CACHE_EXTENSION_INIT_STARTUP"
	SecretPaths = "CACHE_EXTENSION_SECRET_PATHS"
	SecretsRegion = "CACHE_EXTENSION_SECRETS_REGION"
)

// Struct for storing CacheConfiguration
type CacheConfig struct {
	Parameters []plugins.ParameterConfiguration
}

var secretNames = os.Getenv(SecretPaths)
var secretsRegion = os.Getenv(SecretsRegion)

var cacheConfig = CacheConfig{}

// Initialize cache and start the background process to refresh cache
func InitCacheExtensions() {
	var secretNamesSlice []string = strings.Split(secretNames, ",")
	var secretNamesWithPathSlice []string
	for _, secretName := range secretNamesSlice {
		secretNamesWithPathSlice = append(secretNamesWithPathSlice, "/aws/reference/secretsmanager/" + secretName)
	}
	cacheConfig.Parameters = append(cacheConfig.Parameters, plugins.ParameterConfiguration{Region: secretsRegion, Names: secretNamesWithPathSlice})

	// Initialize Cache
	InitCache()
	println(plugins.PrintPrefix, "Cache successfully loaded")
}

// Initialize individual cache
func InitCache() {

	// Read Lambda env variable
	var initCache = os.Getenv(InitializeCacheOnStartup)
	var initCacheInBool = false
	if initCache != "" {
		cacheInBool, err := strconv.ParseBool(initCache)
		if err != nil {
			panic(plugins.PrintPrefix + "Error while converting CACHE_EXTENSION_INIT_STARTUP env variable " +
				initCache)
		} else {
			initCacheInBool = cacheInBool
		}
	}

	// Initialize map and load data from individual services if "CACHE_EXTENSION_INIT_STARTUP" = true
	plugins.InitParameters(cacheConfig.Parameters, initCacheInBool)
}

// Route request to corresponding cache handlers
func RouteCache(cacheType string, name string) string {
	switch cacheType {
	case Parameters:
		return plugins.GetParameterCache(name)
	default:
		return ""
	}
}