package stack

func prepareStack(stack, envStack *Stack) error {
	if len(envStack.Version) > 0 {
		stack.Version = envStack.Version
	}

	if len(envStack.Provider.Name) > 0 {
		stack.Provider.Name = envStack.Provider.Name
	}

	if len(envStack.Provider.GatewayURL) > 0 {
		stack.Provider.GatewayURL = envStack.Provider.GatewayURL
	}

	if len(envStack.Hostname) > 0 {
		stack.Hostname = envStack.Hostname
	}

	if len(envStack.Functions) > 0 {
		stack.Functions = prepareFunctions(stack.Functions, envStack.Functions)
	}

	if len(envStack.StackConfig.Environments) > 0 {
		for envName := range envStack.StackConfig.Environments {
			stack.StackConfig.Environments[envName] = envStack.StackConfig.Environments[envName]
		}
	}

	if envStack.StackConfig.Build.Enabled != nil {
		stack.StackConfig.Build.Enabled = envStack.StackConfig.Build.Enabled
	}

	if len(envStack.StackConfig.Build.Args) > 0 {
		for buildArgsName := range envStack.StackConfig.Build.Args {
			stack.StackConfig.Build.Args[buildArgsName] = envStack.StackConfig.Build.Args[buildArgsName]
		}
	}

	if envStack.StackConfig.Build.UseSha {
		stack.StackConfig.Build.UseSha = envStack.StackConfig.Build.UseSha
	}

	if len(envStack.StackConfig.Deploy.Prefix) > 0 {
		stack.StackConfig.Deploy.Prefix = envStack.StackConfig.Deploy.Prefix
	}

	if len(envStack.StackConfig.Deploy.Suffix) > 0 {
		stack.StackConfig.Deploy.Suffix = envStack.StackConfig.Deploy.Suffix
	}

	if len(envStack.StackConfig.Scale.Min) > 0 {
		stack.StackConfig.Scale.Min = envStack.StackConfig.Scale.Min
	}

	if len(envStack.StackConfig.Scale.Max) > 0 {
		stack.StackConfig.Scale.Max = envStack.StackConfig.Scale.Max
	}

	if len(envStack.StackConfig.Requests.Cpu) > 0 {
		stack.StackConfig.Requests.Cpu = envStack.StackConfig.Requests.Cpu
	}

	if len(envStack.StackConfig.Requests.Memory) > 0 {
		stack.StackConfig.Requests.Memory = envStack.StackConfig.Requests.Memory
	}

	if len(envStack.StackConfig.Limits.Cpu) > 0 {
		stack.StackConfig.Limits.Cpu = envStack.StackConfig.Limits.Cpu
	}

	if len(envStack.StackConfig.Limits.Memory) > 0 {
		stack.StackConfig.Limits.Memory = envStack.StackConfig.Limits.Memory
	}

	if len(envStack.Custom) > 0 {
		stack.Custom = envStack.Custom
	}

	if len(envStack.Swagger.File) > 0 {
		stack.Swagger.File = envStack.Swagger.File
	}

	return nil
}

func prepareFunctions(stackFunction, envStackFunction map[string]Function) (functions map[string]Function) {
	functions = make(map[string]Function)
	for functionName := range envStackFunction {
		if _, exists := stackFunction[functionName]; exists {
			var function Function

			function.Name = stackFunction[functionName].Name

			function.Lang = stackFunction[functionName].Lang
			if len(envStackFunction[functionName].Lang) != 0 {
				function.Lang = envStackFunction[functionName].Lang
			}

			function.Handler = stackFunction[functionName].Handler
			if len(envStackFunction[functionName].Handler) != 0 {
				function.Handler = envStackFunction[functionName].Handler
			}

			function.Image = stackFunction[functionName].Image
			if len(envStackFunction[functionName].Image) != 0 {
				function.Image = envStackFunction[functionName].Image
			}

			function.Path = stackFunction[functionName].Path
			if len(envStackFunction[functionName].Path) != 0 {
				function.Path = envStackFunction[functionName].Path
			}

			function.FunctionConfig = stackFunction[functionName].FunctionConfig

			if len(envStackFunction[functionName].FunctionConfig.Environments) != 0 {
				for envName := range envStackFunction[functionName].FunctionConfig.Environments {
					function.FunctionConfig.Environments[envName] = envStackFunction[functionName].FunctionConfig.Environments[envName]
				}
			}

			if envStackFunction[functionName].FunctionConfig.Build.Enabled != nil {
				function.FunctionConfig.Build.Enabled = envStackFunction[functionName].FunctionConfig.Build.Enabled
			}

			if len(envStackFunction[functionName].FunctionConfig.Build.Args) != 0 {
				for buildArgsName := range envStackFunction[functionName].FunctionConfig.Build.Args {
					function.FunctionConfig.Build.Args[buildArgsName] = envStackFunction[functionName].FunctionConfig.Build.Args[buildArgsName]
				}
			}

			if envStackFunction[functionName].FunctionConfig.Build.UseSha {
				function.FunctionConfig.Build.UseSha = envStackFunction[functionName].FunctionConfig.Build.UseSha
			}

			if len(envStackFunction[functionName].FunctionConfig.Deploy.Prefix) != 0 {
				function.FunctionConfig.Deploy.Prefix = envStackFunction[functionName].FunctionConfig.Deploy.Prefix
			}

			if len(envStackFunction[functionName].FunctionConfig.Deploy.Suffix) != 0 {
				function.FunctionConfig.Deploy.Suffix = envStackFunction[functionName].FunctionConfig.Deploy.Suffix
			}

			if len(envStackFunction[functionName].FunctionConfig.Requests.Memory) != 0 {
				function.FunctionConfig.Requests.Memory = envStackFunction[functionName].FunctionConfig.Requests.Memory
			}

			if len(envStackFunction[functionName].FunctionConfig.Requests.Cpu) != 0 {
				function.FunctionConfig.Requests.Cpu = envStackFunction[functionName].FunctionConfig.Requests.Cpu
			}

			if len(envStackFunction[functionName].FunctionConfig.Limits.Memory) != 0 {
				function.FunctionConfig.Limits.Memory = envStackFunction[functionName].FunctionConfig.Limits.Memory
			}

			if len(envStackFunction[functionName].FunctionConfig.Limits.Cpu) != 0 {
				function.FunctionConfig.Limits.Cpu = envStackFunction[functionName].FunctionConfig.Limits.Cpu
			}

			if len(envStackFunction[functionName].FunctionConfig.Scale.Min) != 0 {
				function.FunctionConfig.Scale.Min = envStackFunction[functionName].FunctionConfig.Scale.Min
			}

			if len(envStackFunction[functionName].FunctionConfig.Scale.Max) != 0 {
				function.FunctionConfig.Scale.Max = envStackFunction[functionName].FunctionConfig.Scale.Max
			}

			if len(envStackFunction[functionName].Plugins) == 0 {
				function.Plugins = stackFunction[functionName].Plugins
			} else {
				function.Plugins = stackFunction[functionName].Plugins
				for pluginName := range envStackFunction[functionName].Plugins {
					if _, exists := stackFunction[functionName].Plugins[pluginName]; exists {
						var plugin Plugin
						if len(envStackFunction[functionName].Plugins[pluginName].Name) == 0 {
							plugin.Name = stackFunction[functionName].Plugins[pluginName].Name
						} else {
							plugin.Name = envStackFunction[functionName].Plugins[pluginName].Name
						}

						if len(envStackFunction[functionName].Plugins[pluginName].Type) == 0 {
							plugin.Type = stackFunction[functionName].Plugins[pluginName].Type
						} else {
							plugin.Type = envStackFunction[functionName].Plugins[pluginName].Type
						}

						if len(envStackFunction[functionName].Plugins[pluginName].Global) == 0 {
							plugin.Global = stackFunction[functionName].Plugins[pluginName].Global
						} else {
							plugin.Global = envStackFunction[functionName].Plugins[pluginName].Global
						}

						if len(envStackFunction[functionName].Plugins[pluginName].Config) == 0 {
							plugin.Config = stackFunction[functionName].Plugins[pluginName].Config
						} else {
							plugin.Config = envStackFunction[functionName].Plugins[pluginName].Config
						}

						if len(envStackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Name) == 0 {
							plugin.ConfigFrom.SecretKeyRef.Name = stackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Name
						} else {
							plugin.ConfigFrom.SecretKeyRef.Name = envStackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Name
						}

						if len(envStackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Key) == 0 {
							plugin.ConfigFrom.SecretKeyRef.Key = stackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Key
						} else {
							plugin.ConfigFrom.SecretKeyRef.Key = envStackFunction[functionName].Plugins[pluginName].ConfigFrom.SecretKeyRef.Key
						}

						function.Plugins[pluginName] = plugin
					} else {
						function.Plugins[pluginName] = envStackFunction[functionName].Plugins[pluginName]
					}
				}
			}
			functions[functionName] = function
		} else {
			functions[functionName] = envStackFunction[functionName]
		}
	}

	for functionName := range stackFunction {
		if _, exists := envStackFunction[functionName]; !exists {
			functions[functionName] = stackFunction[functionName]
		}
	}

	return functions
}
