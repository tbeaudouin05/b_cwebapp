import React, { Component } from 'react';

export default class BmlDgkTableRow extends Component {

    constructor(props) {
        super(props);

        this.input = React.createRef();
    }

    onClick(e) {
        e.preventDefault();
        const { ApplyManualMatching } = this.props;
        ApplyManualMatching(this.props.DgkRowValue.DgkIDCatalogConfig, this.props.DgkRowValue.DgkImgLink);
    }
    //
    render() {
        return (
            <tr key={this.props.rowKey}>
                <td key='3'>
                    <div className="row">
                        <div className="col-xs-4">
                            <div className="img__wrap">
                                <img src={this.props.DgkRowValue.DgkImgLink} alt={this.props.DgkRowValue.DgkSKUName} className="img-thumbnail" height="85" width="85" />
                                <p className="img__description">{this.props.DgkRowValue.DgkSKUName}</p>
                            </div>
                        </div>
                        <td key='2'>
                        <div className="col-xs-3">
                            <form className="row_button" action={this.props.DgkRowValue.DgkSKULink} target="_blank">
                                <button type="submit" className="btn btn-primary mb-2" >See on Digikala</button>
                            </form>
                        </div> </td>
                        <div className="col-xs-3">
                            <form className="row_button" onClick={this.onClick.bind(this)} /*action={ this.handleMatching(this.props.BmlID, this.props.DgkRowValue.DgklIDCatalogConfig)}*/ target="_blank">
                                <button type="submit" className="btn btn-primary mb-2" >Match Manually</button>
                            </form>
                        </div>

                    </div>

                </td>
                <td key='4'> {this.props.DgkRowValue.DgkSKUPrice} </td>

            </tr>
        )
    }
}
