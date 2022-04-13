package handler

import (
	sender_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_blobs "getsturdy.com/api/pkg/blobs/service"
	db_changes "getsturdy.com/api/pkg/changes/db"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/di"
	routes_file "getsturdy.com/api/pkg/file/routes"
	worker_gc "getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/graphql"
	"getsturdy.com/api/pkg/logger"
	db_mutagen "getsturdy.com/api/pkg/mutagen/db"
	service_notifications "getsturdy.com/api/pkg/notification/service"
	db_pki "getsturdy.com/api/pkg/pki/db"
	service_presence "getsturdy.com/api/pkg/presence/service"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	service_sync "getsturdy.com/api/pkg/sync/service"
	uploader_avatars "getsturdy.com/api/pkg/users/avatars/uploader"
	db_users "getsturdy.com/api/pkg/users/db"
	service_users "getsturdy.com/api/pkg/users/service/module"
	db_views "getsturdy.com/api/pkg/view/db"
	meta_view "getsturdy.com/api/pkg/view/meta"
	service_view "getsturdy.com/api/pkg/view/service"
	db_waitinglist "getsturdy.com/api/pkg/waitinglist"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_users.Module)
	c.Import(service_analytics.Module)
	c.Import(db_waitinglist.Module)
	c.Import(db_codebases.Module)
	c.Import(db_views.Module)
	c.Import(db_workspaces.Module)
	c.Import(db_pki.Module)
	c.Import(worker_snapshots.Module)
	c.Import(db_changes.Module)
	c.Import(worker_gc.Module)
	c.Import(service_comments.Module)
	c.Import(service_activity.Module)
	c.Import(service_view.Module)
	c.Import(service_users.Module)
	c.Import(sender_activity.Module)
	c.Import(workers_ci.Module)
	c.Import(service_notifications.Module)
	c.Import(service_workspaces.Module)
	c.Import(meta_view.Module)
	c.Import(db_mutagen.Module)
	c.Import(service_presence.Module)
	c.Import(service_sync.Module)
	c.Import(service_codebases.Module)
	c.Import(service_auth.Module)
	c.Import(service_blobs.Module)
	c.Import(uploader_avatars.Module)
	c.Import(routes_file.Module)
	c.Import(graphql.Module)

	c.Register(ProvideHandler)
}
