CREATE TABLE IF NOT EXISTS IssueRegistry (
  IssueID SERIAL NOT NULL,
  ISBN VARCHAR(17) NOT NULL,
  ReaderID INT NOT NULL,
  IssueApproverID VARCHAR,
  IssueStatus VARCHAR,
  IssueDate DATE,
  ExpectedReturnDate DATE,
  ReturnDate DATE,
  ReturnApproverID INT,
  PRIMARY KEY (IssueID),
  FORIEGN KEY (ISBN) REFERENCES BookInventory (ISBN),
  FORIEGN KEY (ReaderID) REFERENCES Users (ID),
  FORIEGN KEY (IssueApproverID) REFERENCES Users (ID),
  FORIEGN KEY (ReturnApproverID) REFERENCES Users (ID)
);
