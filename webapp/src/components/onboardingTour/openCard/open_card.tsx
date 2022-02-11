// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'

import {bottom} from '@popperjs/core'

import {FormattedMessage} from 'react-intl'

import TourTip from '../../tutorial_tour_tip/tutorial_tour_tip'
import {useMeasurePunchouts} from '../../tutorial_tour_tip/hooks'
import {useAppSelector} from '../../../store/hooks'
import {getOnboardingTourStep} from '../../../store/users'

import {TOUR_BASE} from '../index'

import './open_card.scss'
import {OnboardingCardClassName} from '../../kanban/kanbanCard'

type Props = {
    onPunchholeClick: (e: React.MouseEvent) => void
}

const OpenCardTourStep = (props: Props): JSX.Element => {
    const title = (
        <FormattedMessage
            id='OnboardingTour.OpenACard.Title'
            defaultMessage='Open a card'
        />
    )
    const screen = (
        <FormattedMessage
            id='OnboardingTour.OpenACard.Body'
            defaultMessage='Open a card to explore the powerful ways that Boards can help you organize your work.'
        />
    )

    const punchout = useMeasurePunchouts([`.${OnboardingCardClassName}`], [])
    const currentStep = parseInt(useAppSelector(getOnboardingTourStep), 10)

    return (
        <TourTip
            screen={screen}
            title={title}
            punchOut={punchout}
            step={currentStep}
            tutorialCategory={TOUR_BASE}
            autoTour={true}
            placement={bottom}
            className='OpenCardTourStep'
            hideBackdrop={false}
            clickThroughPunchhole={true}
            hideNavButtons={true}
            singleTip={true}
            onPunchholeClick={props.onPunchholeClick}
            skipCategoryFromBackdrop={true}
            telemetryTag='tourPoint1'
        />
    )
}

export default OpenCardTourStep