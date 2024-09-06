Feature: Lightning Stream, baseline
  Scenario: Get records put in source LMDB environment, from target environment
    Given there is a new LMDB environment "source" with 2 DBs at most
    And there is a new LMDB environment "target" with 2 DBs at most
    And there is a new Minio server
    And there is a bucket "bucket"
    And there is an LS instance syncing LMDB env. "source" to bucket "bucket"
    And there is an LS instance syncing LMDB env. "target" to bucket "bucket"
    When I put an LS-native record "a" "A" in DB "upper" of LMDB env. "source"
    And I put an LS-native record "B" "b" in DB "lower" of LMDB env. "source"
    Then I should get LS-native record "a" "A" in DB "upper" of env. "target"
    And I should get LS-native record "B" "b" in DB "lower" of env. "target"
