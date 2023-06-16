import React, { Component } from "react";

export default class BestSolutionsList extends Component {
    render() {
        const headerStyle = {
            padding: "10px",
            backgroundColor: "#DADBDD",
            fontFamily: "monospace",
            fontSize: 15
        }

        const dataNormalStyle = {
            paddingRight: "10px",
            paddingLeft: "10px",
            fontFamily: "monospace",
            fontWeight: 1,
            fontSize: 14
        }


        if (typeof this.props.solutions === "undefined") {
            return (<div></div>)
        }

        return (
            <div>
                <h3>Best solutions table:</h3>
                <table border="1px" style={{ "borderCollapse": "collapse", width: "900px", padding: "10px" }}>
                    <thead >
                        <tr>
                            <th style={headerStyle}>Number</th>
                            <th style={headerStyle}>FoundedAt</th>
                            <th style={headerStyle}>Sequence</th>
                            <th style={headerStyle}>Deviation</th>
                            <th style={headerStyle}>Text</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            this.props.solutions.map((solution, index) => (
                                <tr key={index}>
                                    <td style={dataNormalStyle}>{solution.Number}</td>
                                    <td style={dataNormalStyle}>{solution.FoundedAt}</td>
                                    <td style={dataNormalStyle}>{solution.Sequence}</td>
                                    <td style={dataNormalStyle}>{solution.Deviation}</td>
                                    <td style={dataNormalStyle}>{solution.Text}</td>
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>
        );
    }
}