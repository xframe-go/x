package event

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Plugin struct {
	eventBus *Bus[interface{}]
	config   GormPluginConfig
}

type GormPluginConfig struct {
	PublishCreated bool
	PublishUpdated bool
	PublishDeleted bool
	Prefix         string
}

func NewPlugin(bus *Bus[any], config GormPluginConfig) *Plugin {
	return &Plugin{
		eventBus: bus,
		config:   config,
	}
}

func (p *Plugin) Name() string {
	return "event-bus-plugin"
}

func (p *Plugin) Initialize(db *gorm.DB) error {
	callback := db.Callback()

	if p.config.PublishCreated {
		err := callback.Create().After("event-bus-publish-created").Register("publish-created", p.publishCreated)
		if err != nil {
			return err
		}
	}

	if p.config.PublishUpdated {
		err := callback.Update().After("event-bus-publish-updated").Register("publish-updated", p.publishUpdated)
		if err != nil {
			return err
		}
	}

	if p.config.PublishDeleted {
		err := callback.Delete().After("event-bus-publish-deleted").Register("publish-deleted", p.publishDeleted)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) publishCreated(db *gorm.DB) {
	model := p.extractModel(db)
	if model == nil {
		return
	}

	eventType := p.getEventType(model, "created")
	p.eventBus.Publish(eventType, map[string]interface{}{
		"action": "created",
		"model":  model,
	})
}

func (p *Plugin) publishUpdated(db *gorm.DB) {
	newModel := p.extractModel(db)
	if newModel == nil {
		return
	}

	oldModel := db.Statement.ReflectValue

	eventType := p.getEventType(newModel, "updated")
	p.eventBus.Publish(eventType, map[string]interface{}{
		"action":   "updated",
		"model":    newModel,
		"oldModel": oldModel,
	})
}

func (p *Plugin) publishDeleted(db *gorm.DB) {
	model := p.extractModel(db)
	if model == nil {
		return
	}

	eventType := p.getEventType(model, "deleted")
	p.eventBus.Publish(eventType, map[string]interface{}{
		"action": "deleted",
		"model":  model,
	})
}

func (p *Plugin) extractModel(db *gorm.DB) interface{} {
	if db.Statement.Dest != nil {
		return db.Statement.Dest
	}

	if db.Statement.Model != nil {
		return db.Statement.Model
	}

	if db.Statement.ReflectValue.IsValid() && !db.Statement.ReflectValue.IsZero() {
		return db.Statement.ReflectValue.Interface()
	}

	return nil
}

func (p *Plugin) getEventType(model interface{}, action string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	typeName := modelType.Name()
	if p.config.Prefix != "" {
		typeName = p.config.Prefix + "." + typeName
	}

	return strings.ToLower(typeName + "." + action)
}
