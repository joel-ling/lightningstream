Feature: Lightning Stream, baseline
  Scenario: Get records put in source LMDB environment, from target environment
    Given there is a new LMDB environment "source.lmdb" with 2 DBs at most
    And there is a new LMDB environment "target.lmdb" with 2 DBs at most
    And there is a new Minio server
    And there is a new bucket "bucket"
    And there is a new LS instance syncing "source.lmdb" to "bucket"
    And there is a new LS instance syncing "target.lmdb" to "bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "A" in DB "upper"
    And in the transaction I put an LS-native record "B" "b" in DB "lower"
    And I commit the transaction
    Then I should get LS-native record "a" "A" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "b" in DB "lower" of "target.lmdb"

  Scenario: Get records put in one LMDB environment from another, both ways
    Given there is a new LMDB environment "active-0.lmdb" with 2 DBs at most
    And there is a new LMDB environment "active-1.lmdb" with 2 DBs at most
    And there is a new Minio server
    And there is a new bucket "bucket"
    And there is a new LS instance syncing "active-0.lmdb" to "bucket"
    And there is a new LS instance syncing "active-1.lmdb" to "bucket"
    When I begin a transaction in "active-0.lmdb"
    And in the transaction I put an LS-native record "a" "A" in DB "upper"
    And I commit the transaction
    And I begin a transaction in "active-1.lmdb"
    And in the transaction I put an LS-native record "B" "b" in DB "lower"
    And I commit the transaction
    Then I should get LS-native record "a" "A" in DB "upper" of "active-1.lmdb"
    And I should get LS-native record "B" "b" in DB "lower" of "active-0.lmdb"

  Scenario: Simulate in-order object replication between S3 buckets
    Given there is a new LMDB environment "source.lmdb" with 2 DBs at most
    And there is a new LMDB environment "target.lmdb" with 2 DBs at most
    And there is a new Minio server
    And there is a new bucket "source.bucket"
    And there is a new bucket "target.bucket"
    And there is a new LS instance syncing "source.lmdb" to "source.bucket"
    And there is a new LS instance syncing "target.lmdb" to "target.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "A" in DB "upper"
    And in the transaction I put an LS-native record "B" "b" in DB "lower"
    And I commit the transaction
    Then I should count a total of 1 object in "source.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "AA" in DB "upper"
    And in the transaction I put an LS-native record "B" "bb" in DB "lower"
    And I commit the transaction
    Then I should count a total of 2 objects in "source.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "AAA" in DB "upper"
    And in the transaction I put an LS-native record "B" "bbb" in DB "lower"
    And I commit the transaction
    Then I should count a total of 3 objects in "source.bucket"
    When I copy object i=0 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "A" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "b" in DB "lower" of "target.lmdb"
    When I copy object i=1 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "AA" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "bb" in DB "lower" of "target.lmdb"
    When I copy object i=2 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "AAA" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "bbb" in DB "lower" of "target.lmdb"

  Scenario: Simulate out-of-order object replication between S3 buckets
    Given there is a new LMDB environment "source.lmdb" with 2 DBs at most
    And there is a new LMDB environment "target.lmdb" with 2 DBs at most
    And there is a new Minio server
    And there is a new bucket "source.bucket"
    And there is a new bucket "target.bucket"
    And there is a new LS instance syncing "source.lmdb" to "source.bucket"
    And there is a new LS instance syncing "target.lmdb" to "target.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "A" in DB "upper"
    And in the transaction I put an LS-native record "B" "b" in DB "lower"
    And I commit the transaction
    Then I should count a total of 1 object in "source.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "AA" in DB "upper"
    And in the transaction I put an LS-native record "B" "bb" in DB "lower"
    And I commit the transaction
    Then I should count a total of 2 objects in "source.bucket"
    When I begin a transaction in "source.lmdb"
    And in the transaction I put an LS-native record "a" "AAA" in DB "upper"
    And in the transaction I put an LS-native record "B" "bbb" in DB "lower"
    And I commit the transaction
    Then I should count a total of 3 objects in "source.bucket"
    When I copy object i=0 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "A" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "b" in DB "lower" of "target.lmdb"
    When I copy object i=2 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "AAA" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "bbb" in DB "lower" of "target.lmdb"
    When I copy object i=1 from "source.bucket" to "target.bucket"
    Then I should get LS-native record "a" "AAA" in DB "upper" of "target.lmdb"
    And I should get LS-native record "B" "bbb" in DB "lower" of "target.lmdb"

  Scenario: Get record put in one LMDB environment and LS-deleted in another
    Given there is a new LMDB environment "writer.lmdb" with 1 DB at most
    And there is a new LMDB environment "eraser.lmdb" with 1 DB at most
    And there is a new Minio server
    And there is a new bucket "bucket"
    And there is a new LS instance syncing "writer.lmdb" to "bucket"
    And there is a new LS instance syncing "eraser.lmdb" to "bucket"
    When I begin a transaction in "writer.lmdb"
    And in the transaction I put an LS-native record "a" "A" in DB "upper"
    And I commit the transaction
    Then I should get LS-native record "a" "A" in DB "upper" of "eraser.lmdb"
    When I begin a transaction in "eraser.lmdb"
    And in the transaction I mark LS-native record "a" in DB "upper" as deleted
    And I commit the transaction
    Then I should get LS-native record "a" "" in DB "upper" of "writer.lmdb"
