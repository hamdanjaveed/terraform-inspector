package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	rowFormat = lipgloss.NewStyle().
			PaddingLeft(2)
	cursorFormat = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("6"))
	selectedFormat = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingTop(2).
			Foreground(lipgloss.Color("2"))

	bodyStyle = lipgloss.NewStyle().
			Margin(1, 1).
			Foreground(lipgloss.Color("2"))

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

type state int

const (
	listState state = iota
	changeState
)

type Bubble struct {
	OutsideChanges ResourceChanges
	Actions        ResourceChanges
	Summary        string

	state        state
	height       int
	heightMargin int
	width        int
	widthMargin  int
	style        *styles
	boxes        []tea.Model

	Cursor        int
	ShowingDetail *int
	// Selected      map[int]struct{}

	ready    bool
	viewport viewport.Model
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
		tm   tea.Model
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(b.headerView())
		footerHeight := lipgloss.Height(b.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !b.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			b.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight-1)
			b.viewport.YPosition = headerHeight
			b.viewport.HighPerformanceRendering = true
			b.viewport.SetContent(thingy)
			b.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			b.viewport.YPosition = headerHeight + 1
			b.viewport.Style = bodyStyle
		} else {
			b.viewport.Width = msg.Width
			b.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	if b.ShowingDetail == nil {
		tm, cmd = updateOverview(msg, b)
	} else {
		tm, cmd = updateDetail(msg, b)
	}
	cmds = append(cmds, cmd)

	if tt, ok := tm.(Bubble); ok {
		tt.viewport, cmd = b.viewport.Update(msg)
		cmds = append(cmds, cmd)
		return tt, tea.Batch(cmds...)
	}

	return tm, tea.Batch(cmds...)
}

func (b Bubble) headerView() string {
	title := titleStyle.Render(b.Summary)
	line := strings.Repeat("─", max(0, b.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (b Bubble) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", b.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, b.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b Bubble) View() string {
	if !b.ready {
		return "\n  Loading..."
	}

	// TODO: i think i need to only use the viewport on the detail screen, and keep the
	// normal screen as is, so each time we transition screens we init a new viewport

	// var s string
	// if m.ShowingDetail == nil {
	// 	s = viewOverview(m)
	// } else {
	// 	s = viewDetail(m)
	// }

	// m.viewport.SetContent(s)
	b.viewport.SetContent(thingy)

	return fmt.Sprintf("%s\n%s\n%s", b.headerView(), b.viewport.View(), b.footerView())
}

func updateDetail(msg tea.Msg, m Bubble) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "up", "k":
		// if m.Cursor > 0 {
		// 	m.Cursor--
		// }
		case "down", "j":
			// m.viewport.LineDown(1)
		// if m.Cursor < len(m.OutsideChanges)+len(m.Actions)-1 {
		// 	m.Cursor++
		// }
		case "esc", "enter", " ":
			m.ShowingDetail = nil
			// _, ok := m.Selected[m.Cursor]
			// if ok {
			// 	delete(m.Selected, m.Cursor)
			// } else {
			// 	m.Selected[m.Cursor] = struct{}{}
			// }
		}
	}
	return m, nil
}

func updateOverview(msg tea.Msg, m Bubble) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.OutsideChanges)+len(m.Actions)-1 {
				m.Cursor++
			}
		case "enter", " ":
			m.ShowingDetail = &m.Cursor
			// _, ok := m.Selected[m.Cursor]
			// if ok {
			// 	delete(m.Selected, m.Cursor)
			// } else {
			// 	m.Selected[m.Cursor] = struct{}{}
			// }
		}
	}
	return m, nil
}

func viewDetail(m Bubble) string {
	return titleStyle.Render(append(m.OutsideChanges, m.Actions...)[*m.ShowingDetail].Diff)
}

func viewOverview(m Bubble) string {
	s := "Outside changes:\n"

	for i, rc := range m.OutsideChanges {
		s += fmt.Sprintf("%s\n", m.viewRow(i, rc))
	}

	s += "\nActions:\n"

	for i, rc := range m.Actions {
		s += fmt.Sprintf("%s\n", m.viewRow(i+len(m.OutsideChanges), rc))
	}

	s += "\nPress q or ctrl+c to quit.\n"
	return s
}

func (b Bubble) viewRow(i int, rc ResourceChange) string {
	cursor := " "
	if b.Cursor == i {
		cursor = ">"
	}

	// selected := false
	// if _, ok := m.Selected[i]; ok {
	// 	selected = true
	// }

	action := rc.Actions.String()
	if !rc.Actions.IsReplace() {
		action = fmt.Sprintf("  %s", rc.Actions)
	}

	s := fmt.Sprintf("%s %s %s", cursor, action, rc.Address)
	if cursor == ">" {
		s = cursorFormat.Render(s)
	} else {
		s = rowFormat.Render(s)
	}

	// if m.Cursor == i {
	// 	s += selectedFormat.Render(rc.Diff)
	// }

	return s
}

const thingy = `
aws_iam_policy.policy5: Refreshing state... [id=arn:aws:iam::941738800554:policy/test_policy_5]
aws_iam_policy.policy2: Refreshing state... [id=arn:aws:iam::941738800554:policy/test_policy_2]
aws_iam_policy.policy: Refreshing state... [id=arn:aws:iam::941738800554:policy/test_policy]
aws_iam_policy.policy3: Refreshing state... [id=arn:aws:iam::941738800554:policy/test_policy_3]
aws_iam_policy.policy4: Refreshing state... [id=arn:aws:iam::941738800554:policy/test_policy_4]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply":

# aws_iam_policy.policy4 has changed
~ resource "aws_iam_policy" "policy4" {
id          = "arn:aws:iam::941738800554:policy/test_policy_4"
name        = "test_policy_4"
+ tags        = {}
# (6 unchanged attributes hidden)
}

# aws_iam_policy.policy5 has changed
~ resource "aws_iam_policy" "policy5" {
id          = "arn:aws:iam::941738800554:policy/test_policy_5"
name        = "test_policy_5"
+ tags        = {}
# (6 unchanged attributes hidden)
}


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the
following plan may include actions to undo or respond to these changes.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following
symbols:
+ create
~ update in-place
- destroy
-/+ destroy and then create replacement

Terraform will perform the following actions:

# aws_iam_policy.policy2 will be updated in-place
~ resource "aws_iam_policy" "policy2" {
id          = "arn:aws:iam::941738800554:policy/test_policy_2"
name        = "test_policy_2"
~ policy      = jsonencode(
~ {
  ~ Statement = [
	  ~ {
		  ~ Action   = [
				"ec2:Describe*",
			  + "sqs:*",
			]
			# (2 unchanged elements hidden)
		},
	]
	# (1 unchanged element hidden)
}
)
tags        = {}
# (5 unchanged attributes hidden)
}

# aws_iam_policy.policy3 must be replaced
-/+ resource "aws_iam_policy" "policy3" {
~ arn         = "arn:aws:iam::941738800554:policy/test_policy_3" -> (known after apply)
~ id          = "arn:aws:iam::941738800554:policy/test_policy_3" -> (known after apply)
~ name        = "test_policy_3" -> "new_test_policy_3" # forces replacement
~ policy_id   = "ANPA5WRABWWVJGIPTT7MC" -> (known after apply)
- tags        = {} -> null
~ tags_all    = {} -> (known after apply)
# (3 unchanged attributes hidden)
}

# aws_iam_policy.policy5 will be destroyed
# (because aws_iam_policy.policy5 is not in configuration)
- resource "aws_iam_policy" "policy5" {
- arn         = "arn:aws:iam::941738800554:policy/test_policy_5" -> null
- description = "My test policy 5" -> null
- id          = "arn:aws:iam::941738800554:policy/test_policy_5" -> null
- name        = "test_policy_5" -> null
- path        = "/" -> null
- policy      = jsonencode(
{
  - Statement = [
	  - {
		  - Action   = [
			  - "ec2:Describe*",
			]
		  - Effect   = "Allow"
		  - Resource = "*"
		},
	]
  - Version   = "2012-10-17"
}
) -> null
- policy_id   = "ANPA5WRABWWVNZ5LCGO57" -> null
- tags        = {} -> null
- tags_all    = {} -> null
}

# aws_instance.app_server will be created
+ resource "aws_instance" "app_server" {
+ ami                                  = "ami-830c94e3"
+ arn                                  = (known after apply)
+ associate_public_ip_address          = (known after apply)
+ availability_zone                    = (known after apply)
+ cpu_core_count                       = (known after apply)
+ cpu_threads_per_core                 = (known after apply)
+ disable_api_termination              = (known after apply)
+ ebs_optimized                        = (known after apply)
+ get_password_data                    = false
+ host_id                              = (known after apply)
+ id                                   = (known after apply)
+ instance_initiated_shutdown_behavior = (known after apply)
+ instance_state                       = (known after apply)
+ instance_type                        = "t4g.nano"
+ ipv6_address_count                   = (known after apply)
+ ipv6_addresses                       = (known after apply)
+ key_name                             = (known after apply)
+ monitoring                           = (known after apply)
+ outpost_arn                          = (known after apply)
+ password_data                        = (known after apply)
+ placement_group                      = (known after apply)
+ placement_partition_number           = (known after apply)
+ primary_network_interface_id         = (known after apply)
+ private_dns                          = (known after apply)
+ private_ip                           = (known after apply)
+ public_dns                           = (known after apply)
+ public_ip                            = (known after apply)
+ secondary_private_ips                = (known after apply)
+ security_groups                      = (known after apply)
+ source_dest_check                    = true
+ subnet_id                            = (known after apply)
+ tags                                 = {
+ "Name" = "ExampleAppServerInstance"
}
+ tags_all                             = {
+ "Name" = "ExampleAppServerInstance"
}
+ tenancy                              = (known after apply)
+ user_data                            = (known after apply)
+ user_data_base64                     = (known after apply)
+ vpc_security_group_ids               = (known after apply)

+ capacity_reservation_specification {
+ capacity_reservation_preference = (known after apply)

+ capacity_reservation_target {
  + capacity_reservation_id = (known after apply)
}
}

+ ebs_block_device {
+ delete_on_termination = (known after apply)
+ device_name           = (known after apply)
+ encrypted             = (known after apply)
+ iops                  = (known after apply)
+ kms_key_id            = (known after apply)
+ snapshot_id           = (known after apply)
+ tags                  = (known after apply)
+ throughput            = (known after apply)
+ volume_id             = (known after apply)
+ volume_size           = (known after apply)
+ volume_type           = (known after apply)
}

+ enclave_options {
+ enabled = (known after apply)
}

+ ephemeral_block_device {
+ device_name  = (known after apply)
+ no_device    = (known after apply)
+ virtual_name = (known after apply)
}

+ metadata_options {
+ http_endpoint               = (known after apply)
+ http_put_response_hop_limit = (known after apply)
+ http_tokens                 = (known after apply)
+ instance_metadata_tags      = (known after apply)
}

+ network_interface {
+ delete_on_termination = (known after apply)
+ device_index          = (known after apply)
+ network_interface_id  = (known after apply)
}

+ root_block_device {
+ delete_on_termination = (known after apply)
+ device_name           = (known after apply)
+ encrypted             = (known after apply)
+ iops                  = (known after apply)
+ kms_key_id            = (known after apply)
+ tags                  = (known after apply)
+ throughput            = (known after apply)
+ volume_id             = (known after apply)
+ volume_size           = (known after apply)
+ volume_type           = (known after apply)
}
}

# aws_sqs_queue.app_queue will be created
+ resource "aws_sqs_queue" "app_queue" {
+ arn                               = (known after apply)
+ content_based_deduplication       = false
+ deduplication_scope               = (known after apply)
+ delay_seconds                     = 0
+ fifo_queue                        = false
+ fifo_throughput_limit             = (known after apply)
+ id                                = (known after apply)
+ kms_data_key_reuse_period_seconds = (known after apply)
+ max_message_size                  = 262144
+ message_retention_seconds         = 345600
+ name                              = "example-queue"
+ name_prefix                       = (known after apply)
+ policy                            = (known after apply)
+ receive_wait_time_seconds         = 0
+ tags                              = {
+ "Name" = "ExampleAppServerInstance"
}
+ tags_all                          = {
+ "Name" = "ExampleAppServerInstance"
}
+ url                               = (known after apply)
+ visibility_timeout_seconds        = 30
}

Plan: 3 to add, 1 to change, 2 to destroy.
`
