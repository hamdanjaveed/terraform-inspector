package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in                     string
		expectedOutsideChanges ResourceChanges
		expectedActions        ResourceChanges
		summary                string
		err                    error
	}{
		{
			in: "oneofeach.tfplan",
			expectedOutsideChanges: ResourceChanges{
				{
					Address: "aws_iam_policy.policy4",
					Type:    "aws_iam_policy",
					Name:    "policy4",
					Actions: Actions{UpdateAction},
					Diff: `~ resource "aws_iam_policy" "policy4" {
        id          = "arn:aws:iam::941738800554:policy/test_policy_4"
        name        = "test_policy_4"
      + tags        = {}
        # (6 unchanged attributes hidden)
    }`,
				},
				{
					Address: "aws_iam_policy.policy5",
					Type:    "aws_iam_policy",
					Name:    "policy5",
					Actions: Actions{UpdateAction},
					Diff: `~ resource "aws_iam_policy" "policy5" {
        id          = "arn:aws:iam::941738800554:policy/test_policy_5"
        name        = "test_policy_5"
      + tags        = {}
        # (6 unchanged attributes hidden)
    }`,
				},
			},
			expectedActions: ResourceChanges{
				{
					Address: "aws_iam_policy.policy2",
					Type:    "aws_iam_policy",
					Name:    "policy2",
					Actions: Actions{UpdateAction},
					Diff: `~ resource "aws_iam_policy" "policy2" {
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
    }`,
				},
				{
					Address: "aws_iam_policy.policy3",
					Type:    "aws_iam_policy",
					Name:    "policy3",
					Actions: Actions{DeleteAction, CreateAction},
					Diff: `-/+ resource "aws_iam_policy" "policy3" {
      ~ arn         = "arn:aws:iam::941738800554:policy/test_policy_3" -> (known after apply)
      ~ id          = "arn:aws:iam::941738800554:policy/test_policy_3" -> (known after apply)
      ~ name        = "test_policy_3" -> "new_test_policy_3" # forces replacement
      ~ policy_id   = "ANPA5WRABWWVJGIPTT7MC" -> (known after apply)
      - tags        = {} -> null
      ~ tags_all    = {} -> (known after apply)
        # (3 unchanged attributes hidden)
    }`,
				},
				{
					Address: "aws_iam_policy.policy5",
					Type:    "aws_iam_policy",
					Name:    "policy5",
					Actions: Actions{DeleteAction},
					Diff: `- resource "aws_iam_policy" "policy5" {
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
    }`,
				},
				{
					Address: "aws_instance.app_server",
					Type:    "aws_instance",
					Name:    "app_server",
					Actions: Actions{CreateAction},
					Diff: `+ resource "aws_instance" "app_server" {
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
    }`,
				},
				{
					Address: "aws_sqs_queue.app_queue",
					Type:    "aws_sqs_queue",
					Name:    "app_queue",
					Actions: Actions{CreateAction},
					Diff: `+ resource "aws_sqs_queue" "app_queue" {
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
    }`,
				},
			},
			summary: "Plan: 3 to add, 1 to change, 2 to destroy.",
		},
	}

	for _, test := range tests {
		f, err := os.Open(filepath.Join("..", "files", "testdata", test.in))
		assert.NoError(t, err)
		s, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		o, a, summary, err := Parse(string(s))
		if test.err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutsideChanges, o, "outside changes are not equal")
			assert.Equal(t, test.expectedActions, a, "actions are not equal")
			assert.Equal(t, test.summary, summary, "summary is not equal")
		}
	}
}
