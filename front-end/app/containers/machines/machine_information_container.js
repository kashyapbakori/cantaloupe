import { connect } from 'react-redux';
import Machine_information from '../../components/machines/Machine_information_component';
import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import * as Actions from '../../actions/actions';

function mapStateToProps(state) {
	return {
		props: state,
	};
}

const mapDispatchToProps = (dispatch, ownProps) => {
	return {
		actions: bindActionCreators(Actions, dispatch),
		dispatch: dispatch,
	};
}

export default connect(mapStateToProps, mapDispatchToProps)(Machine_information);



