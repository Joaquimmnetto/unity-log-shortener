package parser

func DefaultConfig() Config {
	return Config{
		Preprocessors: Preprocessors{
			RemoveAllMatchingFromLine: []string{
				`\033\[\d+m`,      //bash_color_regex
				`</?color?[^>]+>`, //unity_color_regex
			},
		},
		Matchers: Matchers{
			RemoveLine: []string{
				`^\[Code Coverage\].*$`, // code_coverage_start
				`^\[Performance\].*$`,   // peformance_start
				`^(?<namespace>[\w._\x60<>\/\s,]+):(?<method>[\w._\x60<>\/\s,]+) (?<params>\([\w._\x60<>\/\s,&\[\]]*\))\s?(?:\((?<line>at.*)\))?$`, //match stacktrace line. Jesus this one is terrifying.
				`^\s*\(Filename:.*\)`, //log_filename_line
				//`\s*Start importing.*`,                                      //asset_import_line
				//`\s*(\[Worker\s?\w+\]).+`,                                                                 //asset_worker_line
				`^\s*(\[Worker\s?\w+\])  -> \(artifact id:.+`,                                             //asset_worker_import_finished_line
				`^\s*Done ondemand importing asset:.+`,                                                    //asset_worker_import_finished_done_line
				`^Artifact\(content hash=[^\s]+\) downloaded for.*`,                                       //artifact_download_line
				`^ShaderCacheRemote downloaded [\d.]+ bytes for key '\w+'$`,                               //shader_cache_download_line
				`^\s*-\s+Placed sprite\s+ [^\s]+ in page \(.+\) at \(.+\). Page is now \(.+\). \(.+\).*$`, //atlas_sprite_placement_line
				`^Memory Statistics:$`,                                                                    //memory stats [ALLOC_] are removed in tabulated blocks

			},
			RemoveTabulatedBlocks: map[string]TabulatedBlock{
				"licenseServiceConfig":  {Start: `^\[UnityConnectServicesConfig\] Service configuration:$`, MatchStart: true},
				"cacheBlock":            {Start: `^\s*Querying for cacheable assets in Cache Server:`, MatchStart: true},
				"shaderCompilation":     {Start: `^Compiling shader \".+\" pass .+$`, MatchStart: false},
				"shaderSerialization":   {Start: `^Serialized binary data for shader .+ in .+$`, MatchStart: true},
				"memoryStatsAlloc":      {Start: `^\[ALLOC_.+\].*$`, MatchStart: true},
				"domainReloadProfiling": {Start: `^Domain Reload Profiling:$`, MatchStart: true},
				"configParameters":      {Start: `^Configuration Parameters - Can be set up in boot.config$`, MatchStart: true},
				"pluginPreloadSummary":  {Start: `^Preloading .+ native plugins for Editor in .+ ms.*$`, MatchStart: false},
				"assetRefreshSummary":   {Start: `^Asset Pipeline Refresh: Total: .+ seconds - Initiated by .+$`, MatchStart: false},
			},
			RemoveStartEndBlocks: map[string]StartEndBlock{
				"buildReportAssets": {
					Start: `\s*Used Assets and files from the Resources folder, sorted by uncompressed size:`,
					End:   `^---------+$`,
				},
				"playerSizeStats": {
					Start: `^\s*\*\*\*Player size statistics\*\*\*$`,
					End:   `^\s*Total compressed size [\d.]+.+Total uncompressed size [\d.]+.*$`,
				},
				"openCover": {
					Start:      `^\s*Started OpenCover Session$`,
					MatchStart: false,
					End:        `^\s*Finished OpenCover Session$`,
					MatchEnd:   false,
				},
			},
		},
		Summarizers: Summarizers{
			EnableSceneSummarizer:      true,
			EnableAssetsSummarizer:     true,
			EnableCscWarningsSumarizer: true,
		},
	}
}
