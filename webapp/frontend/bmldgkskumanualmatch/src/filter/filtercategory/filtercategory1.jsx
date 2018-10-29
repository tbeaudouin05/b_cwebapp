import React, { Component } from 'react';
import OptionFilterCategory from './option'


export default class FilterCategory1 extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: 'All' };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        this.props.chooseCategory1(event.target.value);
        event.preventDefault();
    }


    render() {
        return (
            <form>
                <select className="form-control" value={this.state.value} onChange={this.handleChange}>
                    {this.props.optionList1.map(option => {
                        return (
                            <OptionFilterCategory optionValue={option.optionValue} optionText={option.optionText} />
                        )
                    })}
                </select>
            </form>
        )
    }
}

