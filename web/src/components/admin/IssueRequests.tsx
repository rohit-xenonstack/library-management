import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { issueAPI } from "../../utils/api";
import styles from "./Admin.module.scss";

export const IssueRequests = () => {
  const queryClient = useQueryClient();
  const { data: requests, isLoading } = useQuery({
    queryKey: ["issueRequests"],
    queryFn: issueAPI.getRequests,
  });

  const approveMutation = useMutation({
    mutationFn: issueAPI.approveRequest,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["issueRequests"] });
    },
  });

  if (isLoading) return <div>Loading requests...</div>;

  return (
    <div className={styles.requestsContainer}>
      <h2>Issue Requests</h2>
      <div className={styles.requestsList}>
        {requests?.data.map((request: any) => (
          <div key={request.id} className={styles.requestCard}>
            <h3>Request #{request.id}</h3>
            <p>Book: {request.book.title}</p>
            <p>User: {request.user.name}</p>
            <p>Status: {request.status}</p>
            {request.status === "pending" && (
              <button
                onClick={() => approveMutation.mutate(request.id)}
                disabled={approveMutation.isPending}
              >
                {approveMutation.isPending ? "Approving..." : "Approve Request"}
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
