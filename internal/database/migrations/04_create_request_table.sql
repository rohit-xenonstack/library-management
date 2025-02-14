CREATE TABLE IF NOT EXISTS RequestEvents (
  ReqID SERIAL NOT NULL,
  BookID VARCHAR(17),
  ReaderID INT,
  RequestDate DATE,
  ApprovalDate DATE,
  ApproverID INT,
  RequestType VARCHAR(25),
  PRIMARY KEY (ReqID),
  FORIEGN KEY (BookID) REFERENCES BookInventory (ISBN),
  FORIEGN KEY (ReaderID) REFERENCES Users (ID),
  FORIEGN KEY (ApproverID) REFERENCES Users (ID)
);
