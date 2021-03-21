package ecs

type TradeHandler struct{}

func (h *TradeHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TRY_TRADE {
		tryTradeEvent := event.data.(TryTrade)

		// trim offering and requested to what is actually possesed
		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)
		offered := make([]Entity, 0)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)
			for item := range tryTradeEvent.offering {
				if inventoryComponent.Items[item] {
					offered = append(offered, item)
				}
			}
		}

		inventoryData, hasInventory = m.getComponent(tryTradeEvent.who, INVENTORY)
		requested := make([]Entity, 0)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)
			for item := range tryTradeEvent.requesting {
				if inventoryComponent.Items[item] {
					requested = append(requested, item)
				}
			}
		}

		// how do we decide?
		// currently accepts any trade

		// need positions to move items between the two
		offeringPositionData, offeringHasPosition := m.getComponent(event.entity, POSITION)
		requestingPositionData, requestingHasPosition := m.getComponent(tryTradeEvent.who, POSITION)

		if offeringHasPosition && requestingHasPosition {
			offeringPositionComponent := offeringPositionData.(Position)
			requestingPositionComponent := requestingPositionData.(Position)

			dx := requestingPositionComponent.X - offeringPositionComponent.X
			dy := requestingPositionComponent.Y - offeringPositionComponent.Y

			// ISSUE: when trading weapons, the move actions cause the person to get damaged by the weapon
			for _, item := range offered {
				returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{item}, event.entity})
				returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{dx, dy}, item})
				returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{item}, tryTradeEvent.who})
			}
			for _, item := range requested {
				returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{item}, tryTradeEvent.who})
				returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{-dx, -dy}, item})
				returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{item}, event.entity})
			}
		}
	}

	return returnEvents
}
