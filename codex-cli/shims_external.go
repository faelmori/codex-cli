package main

import (
	"package-manager-detector"
	"fast-npm-meta"
	"semver"
)

type AgentName string

const (
	Npm   AgentName = "npm"
	Pnpm  AgentName = "pnpm"
	Yarn  AgentName = "yarn"
	Bun   AgentName = "bun"
	Deno  AgentName = "deno"
)

func GetUserAgent() AgentName {
	return package_manager_detector.GetUserAgent()
}

type LatestVersionMeta struct {
	Version string
}

func GetLatestVersion(pkgName string, opts map[string]interface{}) (LatestVersionMeta, error) {
	meta, err := fast_npm_meta.GetLatestVersion(pkgName, opts)
	if err != nil {
		return LatestVersionMeta{}, err
	}
	return LatestVersionMeta{Version: meta.Version}, nil
}

func SemverGt(v1, v2 string) bool {
	return semver.Gt(v1, v2)
}
