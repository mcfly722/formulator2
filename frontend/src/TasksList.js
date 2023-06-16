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
                //console.log(res.data)
            }
        })
    }

    render() {

        const headerStyle = {
            padding: "10px",
            backgroundColor: "#DADBDD",
            fontFamily: "monospace"
        }

        const dataNormalStyle = {
            padding: "10px",
            fontFamily: "monospace"
        }

        const dataOutdatedStyle = {
            padding: "10px",
            fontFamily: "monospace",
            backgroundColor: "#FF6347",
        }

        const dataDoneStyle = {
            padding: "10px",
            fontFamily: "monospace",
            backgroundColor: "#50C878",
        }

        function taskStyle(task) {
            if (task.Done === true) { return dataDoneStyle }
            if (task.TimeoutedOnSec > 0) { return dataOutdatedStyle }
            return dataNormalStyle
        }

        return (
            <div>
                <Header className="header" as="h2" color="blue">
                    Tasks List:
                </Header>
                <table border="1px" style={{ "borderCollapse": "collapse" }}>
                    <thead >
                        <tr >
                            <th style={headerStyle}>Number</th>
                            <th style={headerStyle}>Sequence</th>
                            <th style={headerStyle}>Agent</th>
                            <th style={headerStyle}>Started At</th>
                            <th style={headerStyle}>Elapsed</th>
                            <th style={headerStyle}>Confirmation</th>
                            <th style={headerStyle}>Timeouted On (Sec)</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            this.state.data.map((task, index) => (
                                <tr key={index}>
                                    <td style={taskStyle(task)}>{task.Number}</td>
                                    <td style={taskStyle(task)}>{task.Sequence}</td>
                                    <td style={taskStyle(task)}>{task.Agent}</td>
                                    <td style={taskStyle(task)}>{task.StartedAt}</td>
                                    <td style={taskStyle(task)}>{task.Elapsed}</td>
                                    <td style={taskStyle(task)}>{task.LastConfirmationAgo}</td>
                                    <td style={taskStyle(task)}>{task.TimeoutedOnSec}</td>
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>

        );
    }
}

export default TasksList;