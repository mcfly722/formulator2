import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon } from "semantic-ui-react";

let endpoint = window.location.href;

class TasksList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            task: "",
            items: [],
        };

    }

    componentDidMount() {
        this.getTasks();
    }

    onChange = (event) => { this.setState({ [event.target.name]: event.target.value }) }

    getTasks = () => {
        axios.get(endpoint + "api/tasks").then((res) => {
            if (res.data) {
                console.log(res.data)
                /*
                this.setState({
                    items: res.data.map((item)=>{
                        let color = "yellow";
                    })
                })
                */
            }
        })
    }

    render() {
        return (
            <div>
                <div className="row">
                    <Header className="header" as="h2" color="blue">
                        ToDoList
                    </Header>
                </div>
                <button onClick={this.getTasks}>
                    Update
                </button>
            </div>
        );
    }
}

export default TasksList;