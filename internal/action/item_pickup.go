package action

import (
	"fmt"
	"github.com/hectorgimenez/koolo/internal/action/step"
	"github.com/hectorgimenez/koolo/internal/config"
	"github.com/hectorgimenez/koolo/internal/game"
	"strings"
)

func (b Builder) ItemPickup() *StaticAction {
	return BuildStatic(func(data game.Data) (steps []step.Step) {
		itemsToPickup := b.getItemsToPickup(data)
		for _, item := range itemsToPickup {
			b.logger.Debug(fmt.Sprintf("Item Detected: %s [%s] at X:%d Y:%d", item.Name, item.Quality, item.Position.X, item.Position.Y))
			steps = append(steps, step.PickupItem(b.logger, item))
		}

		return
	}, CanBeSkipped())
}

func (b Builder) getItemsToPickup(data game.Data) []game.Item {
	missingHealingPotions := b.bm.GetMissingCount(data, game.HealingPotion)
	missingManaPotions := b.bm.GetMissingCount(data, game.ManaPotion)
	missingRejuvenationPotions := b.bm.GetMissingCount(data, game.RejuvenationPotion)
	var itemsToPickup []game.Item
	for _, item := range data.Items.Ground {
		// Pickup potions only if they are required
		if item.IsHealingPotion() {
			if missingHealingPotions == 0 {
				continue
			}
			if b.shouldBePickedUp(data, item) {
				itemsToPickup = append(itemsToPickup, item)
				missingHealingPotions--
			}
			continue
		}
		if item.IsManaPotion() {
			if missingManaPotions == 0 {
				continue
			}
			if b.shouldBePickedUp(data, item) {
				itemsToPickup = append(itemsToPickup, item)
				missingManaPotions--
			}
			continue
		}
		if item.IsRejuvPotion() {
			if missingRejuvenationPotions == 0 {
				continue
			}
			if b.shouldBePickedUp(data, item) {
				itemsToPickup = append(itemsToPickup, item)
				missingRejuvenationPotions--
			}
			continue
		}

		if b.shouldBePickedUp(data, item) {
			itemsToPickup = append(itemsToPickup, item)
			continue
		}
	}

	return itemsToPickup
}

func (b Builder) shouldBePickedUp(d game.Data, i game.Item) bool {
	// Exclude Gheed if we are already have one
	if i.Name == game.ItemGrandCharm && i.Quality == game.ItemQualityUnique {
		for _, invItem := range d.Items.Inventory {
			if invItem.Name == game.ItemGrandCharm && invItem.Quality == game.ItemQualityUnique {
				b.logger.Warn("Gheed's Fortune dropped, but you already have one in the inventory, skipping.")
				return false
			}
		}
	}

	// Check if we should pickup gold, based on amount
	if config.Pickit.PickupGold && strings.EqualFold(i.Name, "Gold") {
		if i.Stats[game.StatGold] >= config.Pickit.MinimumGoldToPickup && d.PlayerUnit.Stats[game.StatGold] < d.PlayerUnit.MaxGold() {
			return true
		}
	}

	return i.PickupPass(false)
}
