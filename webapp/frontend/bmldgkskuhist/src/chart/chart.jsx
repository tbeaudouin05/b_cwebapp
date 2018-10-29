import React, { Component } from 'react';
import ChartistGraph from 'react-chartist';
//import ChartistLegend from 'chartist-plugin-legend';


export default class BmlDgkChart extends Component {

    render() {
        /*let plugins = [
            ChartistLegend({
                legendNames: this.props.LegendName,
                position: this.props.LegendDiv
            })
        ]*/

        // do not change plugins to plugin! This is part of external package!
        return (
            <ChartistGraph data={this.props.Data} type={'Line'} /*options={{plugins}}*/ />
        )
    }
}