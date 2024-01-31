package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.katupy.io/ugen"
)

var mainCmd = &cobra.Command{
	Use:          "ugen",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdoutStat, err := os.Stdout.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat stdout: %w", err)
		}

		ulidAsUuid := viper.GetBool("ulid-as-uuid")

		g := &ugen.Generator{
			AnyCharacter: viper.GetBool("any"),
			Digit:        viper.GetBool("digit"),
			Interval:     viper.GetString("interval"),
			Base64:       viper.GetBool("base64"),
			Hex:          viper.GetBool("hex"),
			Ulid:         viper.GetBool("ulid") || ulidAsUuid,
			UlidAsUuid:   ulidAsUuid,
			Uuid:         viper.GetBool("uuid"),
			Uuid7:        viper.GetBool("uuid7"),
			Lower:        viper.GetBool("lower"),
			Upper:        viper.GetBool("upper"),
			Prefix:       viper.GetString("prefix"),
			Suffix:       viper.GetString("suffix"),
			Separator:    viper.GetString("separator"),
			WithLineFeed: (stdoutStat.Mode() & os.ModeCharDevice) == os.ModeCharDevice,
		}

		if g.Interval != "" {
			if err := g.GenInterval(os.Stdout, viper.GetInt("count")); err != nil {
				return fmt.Errorf("failed to generate interval: %w", err)
			}
		} else {
			if err := g.Gen(os.Stdout, viper.GetInt("count"), viper.GetInt("length")); err != nil {
				return fmt.Errorf("failed to generate: %w", err)
			}
		}

		return nil
	},
}

func init() {
	mainCmd.PersistentFlags().BoolP("any", "a", false, "Use any character from the random generator.")
	if err := viper.BindPFlag("any", mainCmd.PersistentFlags().Lookup("any")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().BoolP("digit", "d", false, "Use only digits. If length > 1, digit will never start with 0.")
	if err := viper.BindPFlag("digit", mainCmd.PersistentFlags().Lookup("digit")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().StringP("interval", "i", "", "Generate a random number within the provided open ended, i.e., [_,_), interval. E.g., -1000,1000. If a single number is provided, the interval begins with zero. Ignores --length.")
	if err := viper.BindPFlag("interval", mainCmd.PersistentFlags().Lookup("interval")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("base64", false, "Output base64 encoded strings. Sets --any.")
	if err := viper.BindPFlag("base64", mainCmd.PersistentFlags().Lookup("base64")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().IntP("count", "c", 1, "The number of unique strings to generate.")
	if err := viper.BindPFlag("count", mainCmd.PersistentFlags().Lookup("count")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("hex", false, "Output hex encoded strings. Sets --any.")
	if err := viper.BindPFlag("hex", mainCmd.PersistentFlags().Lookup("hex")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().IntP("length", "l", 12, "Length of the generated string.")
	if err := viper.BindPFlag("length", mainCmd.PersistentFlags().Lookup("length")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("lower", false, "Output characters in lower case.")
	if err := viper.BindPFlag("lower", mainCmd.PersistentFlags().Lookup("lower")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().String("prefix", "", "Write prefix before each generated string.")
	if err := viper.BindPFlag("prefix", mainCmd.PersistentFlags().Lookup("prefix")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().StringP("separator", "s", "\n", "Separator for generated strings.")
	if err := viper.BindPFlag("separator", mainCmd.PersistentFlags().Lookup("separator")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().String("suffix", "", "Write suffix after each generated string.")
	if err := viper.BindPFlag("suffix", mainCmd.PersistentFlags().Lookup("suffix")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("ulid", false, "Generate a ULID.")
	if err := viper.BindPFlag("ulid", mainCmd.PersistentFlags().Lookup("ulid")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("ulid-as-uuid", false, "Generate a ULID displayed as a UUID. Sets --ulid.")
	if err := viper.BindPFlag("ulid-as-uuid", mainCmd.PersistentFlags().Lookup("ulid-as-uuid")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("upper", false, "Output characters in upper case.")
	if err := viper.BindPFlag("upper", mainCmd.PersistentFlags().Lookup("upper")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("uuid", false, "Generate a random (v4) UUID.")
	if err := viper.BindPFlag("uuid", mainCmd.PersistentFlags().Lookup("uuid")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().Bool("uuid7", false, "Generate a v7 UUID.")
	if err := viper.BindPFlag("uuid7", mainCmd.PersistentFlags().Lookup("uuid7")); err != nil {
		panic(err)
	}
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
