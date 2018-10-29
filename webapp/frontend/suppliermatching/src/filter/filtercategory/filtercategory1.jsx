import React, { Component } from 'react';
import OptionFilterCategory from './option'


export default class FilterCategory1 extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: '' };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        this.props.ChooseCategory1(event.target.value);
        event.preventDefault();
    }


    render() {
        return (
            <form>
                    <p >. </p>

                <select className="form-control" value={this.state.value} onChange={this.handleChange}>
                    {this.props.filterCategory1OptionList.map(option => {
                        return (
                            <OptionFilterCategory optionValue={option.optionValue} optionText={option.optionText} />
                        )
                    })}
                </select>
            </form>
        )
    }
}

