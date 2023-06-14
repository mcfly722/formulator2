import React from "react";
import { Container } from "semantic-ui-react";
import TasksList from "./TasksList.js";

function App() {
  return (
    <div>
      <Container>
        <TasksList />
      </Container>
    </div>
  );
}

export default App;
