import React, { Component } from 'react';
import OptionSortFilter from './option.jsx'


export default class SortFilter extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: '' };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        this.props.chooseSorting(event.target.value);
        event.preventDefault();
    }


    render() {
        return (
            <form className="form-inline">
                <div className="form-group row">
                    <label htmlFor="selectSortingOption" className="col-sm-5 col-form-label ">Sort by: </label>
                    <div className="col-sm-7">
                        <select if="selectSortingOption" className=" form-control-sm" value={this.state.value} onChange={this.handleChange}>
                            {this.props.optionList.map(option => {
                                return (
                                    <OptionSortFilter optionValue={option.optionValue} optionText={option.optionText} />
                                )
                            })}
                        </select>
                    </div>
                </div>

            </form>
        )
    }
}

