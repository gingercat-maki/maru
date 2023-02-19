#!/bin/bash
#
# port forward first to localhost:7233

timestamp() {
  date +"%T"
}

echo "Basic Workflows Start: "
st=$(timestamp)
echo "$st"

for i in {1..10}
do
    echo "Welcome $i times"
    for i in {1..100}
    do
        tctl --namespace performance-test wf start --tq temporal-basic --wt basic-workflow --wtt 5 --et 1800 --if ./scenarios/basic-awx.json
        # tctl --namespace performance-test wf start --tq temporal-bench --wt bench-workflow --wtt 5 --et 1800 --if ./scenarios/bench-awx-basic.json
    done
    sleep 1
done

et=$(timestamp)
echo "Basic Workflows End"
echo "from $st to $et"


# Queries
# tctl -ns performance-test workflow query --workflow_id "3226cffc-5ce9-4539-bd95-190aaa1bdae5" --query_type "histogram_csv"

# Tests
# tctl --namespace performance-test wf start --tq temporal-bench --wt bench-workflow --wtt 5 --et 1800 --if ./scenarios/bench-awx-approval.json
# tctl --namespace performance-test wf start --tq temporal-bench --wt bench-workflow --wtt 5 --et 1800 --if ./scenarios/bench-awx-approval.json
