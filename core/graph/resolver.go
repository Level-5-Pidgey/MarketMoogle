/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (resolver.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package graph

import "MarketMoogleAPI/infrastructure/providers/db"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var dbProv = db.NewDbProvider()

type Resolver struct{}