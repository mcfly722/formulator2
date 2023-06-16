import React, { Component } from "react";
import { Container } from "semantic-ui-react";
import State from "./State.js";
import TasksList from "./TasksList.js";
import NewTaskButton from "./NewTaskButton.js";

export default class App extends Component {
  render() {
    return (
      <div>
        <Container>
          <State />
          <br />
          <NewTaskButton />
          <br />
          <TasksList />
        </Container>
      </div>
    );
  }
}
