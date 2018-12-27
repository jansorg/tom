package cmd

//
// func newCreateFrameCommand(context *GoTimeContext, parent *cobra.Command) *cobra.Command {
// 	var start string
// 	var end string
//
// 	var cmd = &cobra.Command{
// 		Use:   "tag name...",
// 		Short: "create a new frame",
// 		Args:  cobra.MinimumNArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			for _, name := range args {
// 				if _, err := context.Store.AddTag(store.Tag{Name: name}); err != nil {
// 					fatal(err)
// 				}
// 			}
// 		},
// 	}
// 	parent.AddCommand(cmd)
//
// 	cmd.Flags().StringP("start", "", &start, "The start time of the frame")
// 	cmd.Flags().StringP("end", "", &end, "The end time of the frame")
// 	cmd.Flags().StringP("duration", "d", &end, "The end time of the frame")
//
// 	return cmd
// }
