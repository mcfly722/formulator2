import React, { Component, useEffect, useState } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon } from "semantic-ui-react";

let endpoint = window.location.href;

class TasksList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            data: []
        };
    }

    componentDidMount() {
        this.interval = setInterval(() => {
            this.GetTasks();
        }, 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }



    onChange = (event) => { this.setState({ [event.target.name]: event.target.value }) }

    GetTasks = () => {
        axios.get(endpoint + "api/tasks").then((res) => {
            if (res.data) {
                this.setState({ data: res.data })
                //console.log()
            }
        })
    }

    render() {
        return (
            <div>
                <div>
                    <Header className="header" as="h2" color="blue">
                        Tasks List:
                    </Header>
                    <table border="1px" style={{ "border-collapse": "collapse" }}>
                        <thead >
                            <tr >
                                <th style={{ padding: "10px" }}>Number</th>
                                <th style={{ padding: "10px" }}>Sequence</th>
                                <th style={{ padding: "10px" }}>Agent</th>
                                <th style={{ padding: "10px" }}>StartedAt</th>
                                <th style={{ padding: "10px" }}>LastConfirmationAt</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                this.state.data.map((task, index) => (
                                    <tr>
                                        <td style={{ padding: "10px" }}>{task.Number}</td>
                                        <td style={{ padding: "10px" }}>{task.Sequence}</td>
                                        <td style={{ padding: "10px" }}>{task.Agent}</td>
                                        <td style={{ padding: "10px" }}>{task.StartedAt}</td>
                                        <td style={{ padding: "10px" }}>{task.LastConfirmationAt}</td>
                                    </tr>
                                ))
                            }
                        </tbody>
                    </table>
                </div>
            </div >
        );
    }
}

export default TasksList;