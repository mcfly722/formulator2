import React from "react";
import { Container } from "semantic-ui-react";
import TasksList from "./TasksList.js";
import NewTaskButton from "./NewTaskButton.js";

function App() {

  return (
    <div>
      <Container>
        <br />
        <NewTaskButton />
        <br />
        <TasksList />
      </Container>
    </div>
  );
}

export default App;
