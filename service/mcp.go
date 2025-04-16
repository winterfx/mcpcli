package service

//func SetServerState(mcpConfig *McpConfig) error {
//	startTime := time.Now()
//	procs, err := process.Processes()
//	if err != nil {
//		return fmt.Errorf("failed to get processes: %w", err)
//	}
//
//	fmt.Printf("formatProcessInfo took: %s\n", time.Since(startTime))
//
//	expectedCommands := formatExpectedCommand(mcpConfig)
//	p := filterProcess(procs, expectedCommands)
//
//	wg := sync.WaitGroup{}
//	for name, server := range mcpConfig.MCPServers {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			expectedCommand, expectedArgs := preTreatment(server.Command, server.Args)
//			fmt.Printf("expectedCommand: %s, expectedArgs: %v\n", expectedCommand, expectedArgs)
//			if isMatch(p, expectedCommand, expectedArgs) {
//				server.Active = true
//				mcpConfig.MCPServers[name] = server
//			}
//		}()
//	}
//	wg.Wait()
//
//	fmt.Printf("SetServerActive took: %s\n", time.Since(startTime))
//	return nil
//}

//func isMatch(p map[string][]string, expectedCommand string, expectedArgs []string) bool {
//	cmdline, ok := p[expectedCommand]
//	if !ok {
//		return false
//	}
//	cmdlineStr := strings.Join(cmdline, " ")
//	for _, arg := range expectedArgs {
//
//		if !strings.Contains(cmdlineStr, arg) {
//			return false
//		}
//	}
//	return true
//}
//
//func formatExpectedCommand(mcpConfig *McpConfig) []string {
//	expectedCommands := make([]string, 0, len(mcpConfig.MCPServers))
//	wg := sync.WaitGroup{}
//	var mu sync.Mutex
//	for _, server := range mcpConfig.MCPServers {
//		wg.Add(1)
//		go func(server McpServer) {
//			defer wg.Done()
//			baseCmd := normalizeCommand(server.Command)
//
//			mu.Lock()
//			expectedCommands = append(expectedCommands, baseCmd)
//			mu.Unlock()
//		}(server)
//	}
//	wg.Wait()
//	return expectedCommands
//}
//
//func filterProcess(procs []*process.Process, expectCmds []string) map[string][]string {
//	processMap := make(map[string][]string)
//	cmds := strings.Join(expectCmds, " ")
//
//	for _, proc := range procs {
//		if proc.Pid < 1000 {
//			continue
//		}
//
//		name, err := proc.Name()
//		if err != nil || name == "" {
//			continue
//		}
//		if !strings.Contains(cmds, name) {
//			continue
//		}
//
//		cmdLine, err := proc.CmdlineSlice()
//		if err != nil || len(cmdLine) < 1 {
//			continue
//		}
//		if _, ok := processMap[name]; ok {
//			processMap[name] = append(processMap[name], cmdLine[1:]...)
//		} else {
//			processMap[name] = cmdLine[1:]
//		}
//	}
//	return processMap
//}
//
//func preTreatment(command string, args []string) (string, []string) {
//	// 处理命令：根据基名替换为 node
//	baseCmd := normalizeCommand(command)
//	// 处理参数：过滤 -y 并保留原始参数格式
//	newArgs := make([]string, 0, len(args))
//	for _, arg := range args {
//		if arg == "-y" { // 精确匹配 -y 参数
//			continue
//		}
//		newArgs = append(newArgs, arg)
//	}
//	return baseCmd, newArgs // 返回处理后的命令和新参数
//}
//
//func normalizeCommand(command string) string {
//	baseCmd := filepath.Base(command)
//	switch baseCmd {
//	case "npx", "npm", "yarn":
//		return "node"
//	default:
//		return baseCmd
//	}
//}
