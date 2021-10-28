package step

import (
	"fmt"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/bitrise-steplib/steps-xcode-test/simulator"
	"github.com/bitrise-steplib/steps-xcode-test/xcodebuild"
)

func printLastLinesOfXcodebuildTestLog(rawXcodebuildOutput string, isRunSuccess bool) {
	const lastLines = "\nLast lines of the build log:"
	if !isRunSuccess {
		log.Errorf(lastLines)
	} else {
		log.Infof(lastLines)
	}

	fmt.Println(stringutil.LastNLines(rawXcodebuildOutput, 20))

	if !isRunSuccess {
		log.Warnf("If you can't find the reason of the error in the log, please check the xcodebuild_test.log.")
	}

	log.Infof(colorstring.Magenta(`
The log file is stored in $BITRISE_DEPLOY_DIR, and its full path
is available in the $BITRISE_XCODEBUILD_TEST_LOG_PATH environment variable.

If you have the Deploy to Bitrise.io step (after this step),
that will attach the file to your build as an artifact!`))
}

func createConfig(input Input, projectPath string, xcodeMajorVersion int, sim simulator.Simulator) Config {
	return Config{
		ProjectPath: projectPath,
		Scheme:      input.Scheme,
		TestPlan:    input.TestPlan,

		SimulatorID:       sim.ID,
		IsSimulatorBooted: sim.Status != simulatorShutdownState,

		XcodeMajorVersion: xcodeMajorVersion,

		TestRepetitionMode:            input.TestRepetitionMode,
		MaximumTestRepetitions:        input.MaximumTestRepetitions,
		RelaunchTestForEachRepetition: input.RelaunchTestsForEachRepetition,
		RetryTestsOnFailure:           input.RetryTestsOnFailure,

		XCConfigContent:    input.XCConfigContent,
		PerformCleanAction: input.PerformCleanAction,
		XcodebuildOptions:  input.XcodebuildOptions,

		LogFormatter:    input.LogFormatter,
		XcprettyOptions: input.XcprettyOptions,

		CacheLevel: input.CacheLevel,

		CollectSimulatorDiagnostics: exportCondition(input.CollectSimulatorDiagnostics),
		HeadlessMode:                input.HeadlessMode,

		DeployDir: input.DeployDir,
	}
}

func createTestParams(cfg Config, xcresultPath, swiftPackagesPath string) xcodebuild.TestRunParams {
	testParams := xcodebuild.TestParams{
		ProjectPath:                    cfg.ProjectPath,
		Scheme:                         cfg.Scheme,
		Destination:                    fmt.Sprintf("id=%s", cfg.SimulatorID),
		TestPlan:                       cfg.TestPlan,
		TestOutputDir:                  xcresultPath,
		TestRepetitionMode:             cfg.TestRepetitionMode,
		MaximumTestRepetitions:         cfg.MaximumTestRepetitions,
		RelaunchTestsForEachRepetition: cfg.RelaunchTestForEachRepetition,
		XCConfigContent:                cfg.XCConfigContent,
		PerformCleanAction:             cfg.PerformCleanAction,
		RetryTestsOnFailure:            cfg.RetryTestsOnFailure,
		AdditionalOptions:              cfg.XcodebuildOptions,
	}

	return xcodebuild.TestRunParams{
		TestParams:                         testParams,
		LogFormatter:                       cfg.LogFormatter,
		XcprettyOptions:                    cfg.XcprettyOptions,
		RetryOnTestRunnerError:             true,
		RetryOnSwiftPackageResolutionError: true,
		SwiftPackagesPath:                  swiftPackagesPath,
		XcodeMajorVersion:                  cfg.XcodeMajorVersion,
	}
}